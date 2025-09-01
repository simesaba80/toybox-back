package dto

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `json:"id" bun:"id,pk"`
	Name                string    `json:"name" bun:"name,notnull"`
	Email               string    `json:"email" bun:"email,notnull unique"`
	PasswordHash        string    `json:"password_hash" bun:"password_hash"`
	DisplayName         string    `json:"display_name" bun:"display_name,notnull"`
	DiscordToken        string    `json:"discord_token" bun:"discord_token"`
	DiscordRefreshToken string    `json:"discord_refresh_token" bun:"discord_refresh_token"`
	DiscordUserID       string    `json:"discord_user_id" bun:"discord_user_id"`
	Profile             string    `json:"profile" bun:"profile"`
	AvatarURL           string    `json:"avatar_url" bun:"avatar_url"`
	TwitterID           string    `json:"twitter_id" bun:"twitter_id"`
	GithubID            string    `json:"github_id" bun:"github_id"`
	CreatedAt           time.Time `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt           time.Time `json:"updated_at" bun:"updated_at,notnull"`
}
