package usecase

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type DiscordUsecase struct {
	discordRepository repository.DiscordRepository
	userRepository    repository.UserRepository
	tokenProvider     TokenProvider
}

func NewDiscordUsecase(discordRepository repository.DiscordRepository, userRepository repository.UserRepository, tokenProvider TokenProvider) *DiscordUsecase {
	return &DiscordUsecase{
		discordRepository: discordRepository,
		userRepository:    userRepository,
		tokenProvider:     tokenProvider,
	}
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

func (uc *DiscordUsecase) AuthenticateUser(ctx context.Context, code string) (string, entity.DiscordUser, error) {
	token, err := uc.discordRepository.GetDiscordToken(ctx, code)
	if err != nil {
		return "", entity.DiscordUser{}, err
	}

	discordUser, err := uc.discordRepository.FetchDiscordUser(ctx, token)
	if err != nil {
		return "", entity.DiscordUser{}, err
	}
	guildIDs, err := uc.discordRepository.GetDiscordGuilds(ctx, token)
	if err != nil {
		return "", entity.DiscordUser{}, err
	}
	allowedGuildIDs, err := uc.discordRepository.GetAllowedDiscordGuilds(ctx)
	if err != nil {
		return "", entity.DiscordUser{}, err
	}
	if !userBelongsToAllowedGuild(guildIDs, allowedGuildIDs) {
		return "", entity.DiscordUser{}, errors.New("ユーザーは許可されたDiscordギルドに所属していません")
	}

	user, err := uc.userRepository.GetUserByDiscordUserID(ctx, discordUser.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			user, err = uc.userRepository.Create(ctx, &entity.User{
				DiscordUserID:       discordUser.ID,
				DiscordToken:        token.AccessToken,
				DiscordRefreshToken: token.RefreshToken,
			})
			if err != nil {
				return "", entity.DiscordUser{}, err
			}
		} else {
			return "", entity.DiscordUser{}, err
		}
	}

	appToken, err := uc.tokenProvider.GenerateToken(user.ID.String())
	if err != nil {
		return "", entity.DiscordUser{}, err
	}

	return appToken, discordUser, nil
}

func userBelongsToAllowedGuild(guildIDs []string, allowedGuildIDs []string) bool {
	for _, guildID := range guildIDs {
		if slices.Contains(allowedGuildIDs, guildID) {
			return true
		}
	}
	return false
}
