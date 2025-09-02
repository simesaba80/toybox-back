package work

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
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
