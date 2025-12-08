package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetAll(ctx context.Context) ([]*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserByDiscordUserID(ctx context.Context, discordUserID string) (*entity.User, error)
}
