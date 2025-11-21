package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID
	Name                string
	Email               string
	PasswordHash        string
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
