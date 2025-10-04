package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
)

type CommentRepository struct {
	db *bun.DB
}

func NewCommentRepository(db *bun.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) FindByWorkID(ctx context.Context, workID uuid.UUID) ([]*entity.Comment, error) {
	var dtoComments []*dto.Comment
	err := r.db.NewSelect().
		Model(&dtoComments).
		Where("work_id = ?", workID).
		Relation("User").
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	entityComments := make([]*entity.Comment, len(dtoComments))
	for i, dtoComment := range dtoComments {
		entityComments[i] = dtoComment.ToCommentEntity()
	}

	return entityComments, nil
}
