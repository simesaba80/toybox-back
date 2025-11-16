package auth

import (
	"context"
	"errors"

	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
)

type AuthRepository struct {
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) GetDiscordClientID(ctx context.Context) (string, error) {
	if config.DISCORD_CLIENT_ID == "" {
		return "", errors.New("DISCORD_CLIENT_ID is not set")
	}
	return config.DISCORD_CLIENT_ID, nil
}

func (r *AuthRepository) GetHostURL(ctx context.Context) (string, error) {
	if config.HOST_URL == "" {
		return "", errors.New("HOST_URL is not set")
	}
	return config.HOST_URL, nil
}
