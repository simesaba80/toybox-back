package repository

import (
	"context"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type TokenRepository interface {
	Create(ctx context.Context, token *entity.Token) (*entity.Token, error)
	CheckRefreshToken(ctx context.Context, refreshToken string) (string, error)
}
