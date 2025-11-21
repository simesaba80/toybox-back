package schema

import (
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type GetDiscordAuthURLResponse struct {
	URL string `json:"url"`
}

type GetDiscordTokenResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	User         DiscordUserResponse `json:"user"`
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

func ToGetDiscordTokenResponse(token entity.DiscordToken, user entity.DiscordUser) GetDiscordTokenResponse {
	return GetDiscordTokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         ToDiscordUserResponse(user),
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
