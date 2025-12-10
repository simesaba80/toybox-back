package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type TokenRepository interface {
	Create(ctx context.Context, token *entity.Token) (*entity.Token, error)
	CheckRefreshToken(ctx context.Context, refreshToken uuid.UUID) (string, error)
	UpdateRefreshToken(ctx context.Context, refreshToken uuid.UUID) (*entity.Token, error)
}
