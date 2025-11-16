package usecase

import (
	"context"

	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type AuthUsecase struct {
	authRepository repository.AuthRepository
}

func NewAuthUsecase(authRepository repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{authRepository: authRepository}
}

func (uc *AuthUsecase) GetDiscordAuthURL(ctx context.Context) (string, error) {
	DISCORD_CLIENT_ID, err := uc.authRepository.GetDiscordClientID(ctx)
	if err != nil {
		return "", err
	}
	HOST_URL, err := uc.authRepository.GetHostURL(ctx)
	if err != nil {
		return "", err
	}
	authURL := "https://discord.com/oauth2/authorize?client_id=" + DISCORD_CLIENT_ID + "&response_type=code&redirect_uri=" + HOST_URL + "/users/auth/callback&scope=identify+email+guilds"
	return authURL, nil
}
