package repository

import (
	"context"

	"github.com/google/uuid"
)

type TagRepository interface {
	ExistAll(ctx context.Context, ids []uuid.UUID) (bool, error)
}
