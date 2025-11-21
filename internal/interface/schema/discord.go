package schema

import (
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type GetDiscordAuthURLResponse struct {
	URL string `json:"url"`
}

type GetDiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DiscordUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

func ToGetDiscordAuthURLResponse(authURL string) GetDiscordAuthURLResponse {
	return GetDiscordAuthURLResponse{
		URL: authURL,
	}
}

func ToGetDiscordTokenResponse(appToken string, refreshToken string) GetDiscordTokenResponse {
	return GetDiscordTokenResponse{
		AccessToken:  appToken,
		RefreshToken: refreshToken,
	}
}

func ToDiscordUserResponse(user entity.DiscordUser) DiscordUserResponse {
	return DiscordUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
}
