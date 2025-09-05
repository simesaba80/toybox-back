package work

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
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

func (r *WorkRepository) GetAll(ctx context.Context) ([]*entity.Work, error) {
	var dtoWorks []*dto.Work
	err := r.db.NewSelect().Model(&dtoWorks).Relation("Assets").Order("created_at DESC").Limit(20).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ToEntities(dtoWorks), nil
}

func (r *WorkRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	var dtoWork dto.Work
	err := r.db.NewSelect().Model(&dtoWork).Relation("Assets").Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ToEntity(&dtoWork), nil
}

func (r *WorkRepository) Create(ctx context.Context, work *entity.Work) (*entity.Work, error) {
	dtoWork := ToDTO(work)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.NewInsert().Model(dtoWork).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create work in transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return work, nil
}
