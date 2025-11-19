package repository

import (
	"context"
)

type DiscordRepository interface {
	GetDiscordClientID(ctx context.Context) (string, error)
	GetHostURL(ctx context.Context) (string, error)
	GetDiscordAuthURL(ctx context.Context) (string, error)
	GetDiscordToken(ctx context.Context, code string) (string, string, error)
}
