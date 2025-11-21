package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
	"golang.org/x/oauth2"
)

type DiscordRepository struct {
	DiscordOAuthConfig *oauth2.Config
}

const discordUserEndpoint = "https://discord.com/api/v10/users/@me"
const discordGuildsEndpoint = "https://discord.com/api/v10/users/@me/guilds"

type discordUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

type discordGuildResponse struct {
	ID string `json:"id"`
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

func (r *DiscordRepository) FetchDiscordUser(ctx context.Context, token entity.DiscordToken) (entity.DiscordUser, error) {
	oauthToken := toOAuth2Token(token)
	client := r.DiscordOAuthConfig.Client(ctx, oauthToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discordUserEndpoint, nil)
	if err != nil {
		return entity.DiscordUser{}, fmt.Errorf("discord user request作成に失敗しました: %w", err)
	}

	response, err := client.Do(req)
	if err != nil {
		return entity.DiscordUser{}, fmt.Errorf("discord user取得リクエストに失敗しました: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return entity.DiscordUser{}, fmt.Errorf("discord user取得でステータスコード%d: %s", response.StatusCode, string(body))
	}

	var payload discordUserResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return entity.DiscordUser{}, fmt.Errorf("discord userレスポンスのデコードに失敗しました: %w", err)
	}

	return entity.DiscordUser{
		ID:       payload.ID,
		Username: payload.Username,
		Avatar:   payload.Avatar,
		Email:    payload.Email,
	}, nil
}

func (r *DiscordRepository) GetDiscordGuilds(ctx context.Context, token entity.DiscordToken) ([]string, error) {
	oauthToken := toOAuth2Token(token)
	client := r.DiscordOAuthConfig.Client(ctx, oauthToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discordGuildsEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("discord guilds request作成に失敗しました: %w", err)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discord guilds取得リクエストに失敗しました: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return nil, fmt.Errorf("discord guilds取得でステータスコード%d: %s", response.StatusCode, string(body))
	}

	var payload []discordGuildResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("discord guildsレスポンスのデコードに失敗しました: %w", err)
	}

	guildIDs := make([]string, len(payload))
	for i, guild := range payload {
		guildIDs[i] = guild.ID
	}
	return guildIDs, nil
}

func (r *DiscordRepository) GetAllowedDiscordGuilds(ctx context.Context) ([]string, error) {
	return config.DISCORD_GUILD_IDS, nil
}

func toOAuth2Token(token entity.DiscordToken) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
}
