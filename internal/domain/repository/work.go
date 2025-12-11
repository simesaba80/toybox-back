package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type WorkRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Work, int, error)
	GetAllPublic(ctx context.Context, limit, offset int) ([]*entity.Work, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error)
	ExistsById(ctx context.Context, id uuid.UUID) (bool, error)
	Create(ctx context.Context, work *entity.Work) (*entity.Work, error)
}
