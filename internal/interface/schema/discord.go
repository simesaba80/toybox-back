package schema

type GetDiscordAuthURLResponse struct {
	URL string `json:"url"`
}

func ToGetDiscordAuthURLResponse(authURL string) GetDiscordAuthURLResponse {
	return GetDiscordAuthURLResponse{
		URL: authURL,
	}
}
