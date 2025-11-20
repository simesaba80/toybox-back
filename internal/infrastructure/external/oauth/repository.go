package oauth

import (
	"context"
	"errors"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
	"golang.org/x/oauth2"
)

type DiscordRepository struct {
	DiscordOAuthConfig *oauth2.Config
}

func NewDiscordRepository() *DiscordRepository {
	return &DiscordRepository{
		DiscordOAuthConfig: &oauth2.Config{
			ClientID:     config.DISCORD_CLIENT_ID,
			ClientSecret: config.DISCORD_CLIENT_SECRET,
			RedirectURL:  config.HOST_URL + "/auth/discord/callback",
			Scopes:       []string{"identify", "email", "guilds"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		},
	}
}

func (r *DiscordRepository) GetDiscordClientID(ctx context.Context) (string, error) {
	if r.DiscordOAuthConfig.ClientID == "" {
		return "", errors.New("discord client ID is not set")
	}
	return r.DiscordOAuthConfig.ClientID, nil
}

func (r *DiscordRepository) GetHostURL(ctx context.Context) (string, error) {
	if r.DiscordOAuthConfig.RedirectURL == "" {
		return "", errors.New("redirect URL is not set")
	}
	return r.DiscordOAuthConfig.RedirectURL, nil
}

func (r *DiscordRepository) GetDiscordAuthURL(ctx context.Context) (string, error) {
	return r.DiscordOAuthConfig.AuthCodeURL("", oauth2.AccessTypeOffline), nil
}

func (r *DiscordRepository) GetDiscordToken(ctx context.Context, code string) (entity.DiscordToken, error) {
	token, err := r.DiscordOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return entity.DiscordToken{}, err
	}
	return entity.DiscordToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		ExpiresIn:    token.ExpiresIn,
		TokenType:    token.TokenType,
	}, nil
}
