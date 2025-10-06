package comment

import (
	"context"
	"fmt"

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

func (r *CommentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	var dtoComment dto.Comment
	err := r.db.NewSelect().
		Model(&dtoComment).
		Where("comment.id = ?", id).
		Relation("User").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return dtoComment.ToCommentEntity(), nil
}

func (r *CommentRepository) Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
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

	dtoComment := dto.ToCommentDTO(comment)

	_, err = tx.NewInsert().Model(dtoComment).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment in transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return dtoComment.ToCommentEntity(), nil
}
