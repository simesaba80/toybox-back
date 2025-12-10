package tag

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TagRepository struct {
	db *bun.DB
}

func NewTagRepository(db *bun.DB) *TagRepository {
	return &TagRepository{
		db: db,
	}
}

func (r *TagRepository) ExistAll(ctx context.Context, ids []uuid.UUID) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}

	count, err := r.db.NewSelect().
		Table("tag").
		Where("id IN (?)", bun.In(ids)).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count == len(ids), nil
}
