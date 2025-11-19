package schema

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

func ToGetDiscordTokenResponse(accessToken, refreshToken string) GetDiscordTokenResponse {
	return GetDiscordTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
