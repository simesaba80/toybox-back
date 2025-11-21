package usecase

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type AuthUsecase struct {
	discordRepository repository.DiscordRepository
	userRepository    repository.UserRepository
	tokenProvider     TokenProvider
	tokenRepository   repository.TokenRepository
}

func NewAuthUsecase(discordRepository repository.DiscordRepository, userRepository repository.UserRepository, tokenProvider TokenProvider, tokenRepository repository.TokenRepository) *AuthUsecase {
	return &AuthUsecase{
		discordRepository: discordRepository,
		userRepository:    userRepository,
		tokenProvider:     tokenProvider,
		tokenRepository:   tokenRepository,
	}
}

func (uc *AuthUsecase) GetDiscordAuthURL(ctx context.Context) (string, error) {
	if _, err := uc.discordRepository.GetDiscordClientID(ctx); err != nil {
		return "", err
	}
	if _, err := uc.discordRepository.GetHostURL(ctx); err != nil {
		return "", err
	}
	return uc.discordRepository.GetDiscordAuthURL(ctx)
}

func (uc *AuthUsecase) AuthenticateUser(ctx context.Context, code string) (string, string, error) {
	token, err := uc.discordRepository.GetDiscordToken(ctx, code)
	if err != nil {
		return "", "", err
	}

	discordUser, err := uc.discordRepository.FetchDiscordUser(ctx, token)
	if err != nil {
		return "", "", err
	}
	guildIDs, err := uc.discordRepository.GetDiscordGuilds(ctx, token)
	if err != nil {
		return "", "", err
	}
	allowedGuildIDs, err := uc.discordRepository.GetAllowedDiscordGuilds(ctx)
	if err != nil {
		return "", "", err
	}
	if !userBelongsToAllowedGuild(guildIDs, allowedGuildIDs) {
		return "", "", errors.New("ユーザーは許可されたDiscordギルドに所属していません")
	}

	user, err := uc.userRepository.GetUserByDiscordUserID(ctx, discordUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user, err = uc.userRepository.Create(ctx, &entity.User{
				DiscordUserID:       discordUser.ID,
				DiscordToken:        token.AccessToken,
				DiscordRefreshToken: token.RefreshToken,
			})
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	}

	appToken, err := uc.tokenProvider.GenerateToken(user.ID.String())
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.tokenRepository.Create(ctx, &entity.Token{
		UserID: user.ID.String(),
	})
	if err != nil {
		return "", "", err
	}

	return appToken, refreshToken.RefreshToken, nil
}

func userBelongsToAllowedGuild(guildIDs []string, allowedGuildIDs []string) bool {
	for _, guildID := range guildIDs {
		if slices.Contains(allowedGuildIDs, guildID) {
			return true
		}
	}
	return false
}

func (uc *AuthUsecase) RegenerateToken(ctx context.Context, refreshToken string) (string, error) {
	userID, err := uc.tokenRepository.CheckRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	appToken, err := uc.tokenProvider.GenerateToken(userID)
	if err != nil {
		return "", err
	}
	return appToken, nil
}
