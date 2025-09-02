package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type WorkRepository interface {
	GetAll(ctx context.Context) ([]*entity.Work, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error)
}
