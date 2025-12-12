package schema

import (
	"github.com/golang-jwt/jwt/v5"
)

type GetDiscordAuthURLResponse struct {
	URL string `json:"url"`
}

type GetDiscordTokenResponse struct {
	AccessToken string `json:"access_token"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type DiscordUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

type RegenerateTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RegenerateTokenResponse struct {
	AppToken    string `json:"app_token"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func ToGetDiscordAuthURLResponse(authURL string) GetDiscordAuthURLResponse {
	return GetDiscordAuthURLResponse{
		URL: authURL,
	}
}

func ToGetDiscordTokenResponse(appToken string, displayName string, avatarURL string) GetDiscordTokenResponse {
	return GetDiscordTokenResponse{
		AccessToken: appToken,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}
}

func ToRegenerateTokenResponse(appToken string) RegenerateTokenResponse {
	return RegenerateTokenResponse{
		AppToken: appToken,
	}
}
