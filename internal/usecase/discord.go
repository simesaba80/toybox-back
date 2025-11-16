package usecase

import (
	"context"

	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type DiscordUsecase struct {
	authRepository repository.DiscordRepository
}

func NewDiscordUsecase(authRepository repository.DiscordRepository) *DiscordUsecase {
	return &DiscordUsecase{authRepository: authRepository}
}

func (uc *DiscordUsecase) GetDiscordAuthURL(ctx context.Context) (string, error) {
	if _, err := uc.authRepository.GetDiscordClientID(ctx); err != nil {
		return "", err
	}
	if _, err := uc.authRepository.GetHostURL(ctx); err != nil {
		return "", err
	}
	return uc.authRepository.GetDiscordAuthURL(ctx)
}
