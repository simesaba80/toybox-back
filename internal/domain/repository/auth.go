package repository

import (
	"context"
)

type AuthRepository interface {
	GetDiscordClientID(ctx context.Context) (string, error)
	GetHostURL(ctx context.Context) (string, error)
}
