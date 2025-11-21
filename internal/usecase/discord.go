package usecase

import (
	"context"
	"errors"
	"slices"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type DiscordUsecase struct {
	discordRepository repository.DiscordRepository
	userRepository    repository.UserRepository
}

func NewDiscordUsecase(discordRepository repository.DiscordRepository, userRepository repository.UserRepository) *DiscordUsecase {
	return &DiscordUsecase{discordRepository: discordRepository, userRepository: userRepository}
}

func (uc *DiscordUsecase) GetDiscordAuthURL(ctx context.Context) (string, error) {
	if _, err := uc.discordRepository.GetDiscordClientID(ctx); err != nil {
		return "", err
	}
	if _, err := uc.discordRepository.GetHostURL(ctx); err != nil {
		return "", err
	}
	return uc.discordRepository.GetDiscordAuthURL(ctx)
}

func (uc *DiscordUsecase) AuthenticateUser(ctx context.Context, code string) (entity.DiscordToken, entity.DiscordUser, error) {
	token, err := uc.discordRepository.GetDiscordToken(ctx, code)
	if err != nil {
		return entity.DiscordToken{}, entity.DiscordUser{}, err
	}

	user, err := uc.discordRepository.FetchDiscordUser(ctx, token)
	if err != nil {
		return entity.DiscordToken{}, entity.DiscordUser{}, err
	}
	guildIDs, err := uc.discordRepository.GetDiscordGuilds(ctx, token)
	if err != nil {
		return entity.DiscordToken{}, entity.DiscordUser{}, err
	}
	allowedGuildIDs, err := uc.discordRepository.GetAllowedDiscordGuilds(ctx)
	if err != nil {
		return entity.DiscordToken{}, entity.DiscordUser{}, err
	}

	for _, guildID := range guildIDs {
		if guildID == "" {
			continue
		}
		if slices.Contains(allowedGuildIDs, guildID) {
			return token, user, nil
		}
	}

	return entity.DiscordToken{}, entity.DiscordUser{}, errors.New("ユーザーは許可されたDiscordギルドに所属していません")
}
