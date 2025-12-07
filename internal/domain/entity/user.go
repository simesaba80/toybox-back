package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID
	Name                string
	Email               string
	DisplayName         string
	DiscordToken        string
	DiscordRefreshToken string
	DiscordUserID       string
	Profile             string
	AvatarURL           string
	TwitterID           string
	GithubID            string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewUser(name string, email string, displayName string, discordToken string, discordRefreshToken string, discordUserID string, profile string, avatarURL string, twitterID string, githubID string) *User {
	return &User{
		ID:                  uuid.New(),
		Name:                name,
		Email:               email,
		DisplayName:         displayName,
		DiscordToken:        discordToken,
		DiscordRefreshToken: discordRefreshToken,
		DiscordUserID:       discordUserID,
		Profile:             profile,
		AvatarURL:           avatarURL,
		TwitterID:           twitterID,
		GithubID:            githubID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}
