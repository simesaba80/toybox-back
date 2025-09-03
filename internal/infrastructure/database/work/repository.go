package work

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
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
	var works []*entity.Work
	err := r.db.NewSelect().Model(&works).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return works, nil
}

func (r *WorkRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	var work entity.Work
	err := r.db.NewSelect().Model(&work).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func (r *WorkRepository) BeginTx(ctx context.Context) (repository.WorkRepositoryTx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &workRepositoryTx{tx: tx}, nil
}

type workRepositoryTx struct {
	tx bun.Tx
}

func (r *workRepositoryTx) Create(ctx context.Context, work *entity.Work) (*entity.Work, error) {
	_, err := r.tx.NewInsert().Model(work).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return work, nil
}

func (r *workRepositoryTx) Commit() error {
	return r.tx.Commit()
}

func (r *workRepositoryTx) Rollback() error {
	return r.tx.Rollback()
}
