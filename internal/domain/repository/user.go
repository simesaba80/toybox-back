package repository

import (
	"context"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetAll(ctx context.Context) ([]*entity.User, error)
}
