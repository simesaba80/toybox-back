package repository

import (
	"context"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type DiscordRepository interface {
	GetDiscordClientID(ctx context.Context) (string, error)
	GetHostURL(ctx context.Context) (string, error)
	GetDiscordAuthURL(ctx context.Context) (string, error)
	GetDiscordToken(ctx context.Context, code string) (entity.DiscordToken, error)
}
