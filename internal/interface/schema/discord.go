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

func ToGetDiscordAuthURLResponse(authURL string) GetDiscordAuthURLResponse {
	return GetDiscordAuthURLResponse{
		URL: authURL,
	}
}

func ToGetDiscordTokenResponse(token entity.DiscordToken) GetDiscordTokenResponse {
	return GetDiscordTokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}
