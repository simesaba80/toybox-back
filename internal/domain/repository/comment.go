package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type CommentRepository interface {
	FindByWorkID(ctx context.Context, workID uuid.UUID) ([]*entity.Comment, error)
}
