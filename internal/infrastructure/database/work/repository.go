package work

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
)

type WorkRepository struct {
	db *bun.DB
}

func NewWorkRepository(db *bun.DB) *WorkRepository {
	return &WorkRepository{
		db: db,
	}
}

func (r *WorkRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Work, int, error) {
	var dtoWorks []*dto.Work

	total, err := r.db.NewSelect().Model(&dtoWorks).Count(ctx)
	if err != nil {
		return nil, 0, domainerrors.ErrFailedToGetAllWorksByLimitAndOffset
	}

	err = r.db.NewSelect().
		Model(&dtoWorks).
		Relation("Assets").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, domainerrors.ErrWorkNotFound
		}
		return nil, 0, domainerrors.ErrFailedToGetAllWorksByLimitAndOffset
	}

	entityWorks := make([]*entity.Work, len(dtoWorks))
	for i, dtoWork := range dtoWorks {
		entityWorks[i] = dtoWork.ToWorkEntity()
	}

	return entityWorks, total, nil
}

func (r *WorkRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	var dtoWork dto.Work
	err := r.db.NewSelect().Model(&dtoWork).Relation("Assets").Where("id = ?", id).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerrors.ErrWorkNotFound
		}
		return nil, domainerrors.ErrFailedToGetWorkById
	}
	return dtoWork.ToWorkEntity(), nil
}

func (r *WorkRepository) ExistsById(ctx context.Context, id uuid.UUID) (bool, error) {
	var dtoWork dto.Work
	exists, err := r.db.NewSelect().
		Model(&dtoWork).
		Where("id = ?", id).
		Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *WorkRepository) Create(ctx context.Context, work *entity.Work) (*entity.Work, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, domainerrors.ErrFailedToBeginTransaction
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	dtoWork := dto.ToWorkDTO(work)

	_, err = tx.NewInsert().Model(dtoWork).Exec(ctx)
	if err != nil {
		return nil, domainerrors.ErrFailedToCreateWork
	}
	thumbnail := &dto.Thumbnail{
		WorkID:  dtoWork.ID,
		AssetID: dtoWork.ThumbnailAssetID,
	}

	_, err = tx.NewInsert().Model(thumbnail).Exec(ctx)
	if err != nil {
		return nil, domainerrors.ErrFailedToCreateThumbnail
	}

	for _, asset := range dtoWork.Assets {
		_, err = tx.NewUpdate().Model(asset).Set("work_id = ?", dtoWork.ID).Where("id = ?", asset.ID).Exec(ctx)
		if err != nil {
			return nil, domainerrors.ErrFailedToCreateAsset
		}
	}

	_, err = tx.NewUpdate().Model(&dto.Asset{}).Set("work_id = ?", dtoWork.ID).Where("id = ?", dtoWork.ThumbnailAssetID).Exec(ctx)
	if err != nil {
		return nil, domainerrors.ErrFailedToCreateAsset
	}

	if len(dtoWork.URLs) > 0 {
		_, err = tx.NewInsert().Model(&dtoWork.URLs).Exec(ctx)
		if err != nil {
			return nil, domainerrors.ErrFailedToCreateURL
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, domainerrors.ErrFailedToCommitTransaction
	}

	return dtoWork.ToWorkEntity(), nil
}
