package entity

import "time"

type DiscordToken struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	ExpiresIn    int64
	TokenType    string
}

type DiscordUser struct {
	ID         string
	Username   string
	AvatarHash string
	Email      string
}
