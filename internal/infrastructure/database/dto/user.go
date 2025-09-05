package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user"`

	ID                  uuid.UUID `bun:"id,pk"`
	Name                string    `bun:"name,notnull"`
	Email               string    `bun:"email,notnull,unique"`
	PasswordHash        string    `bun:"password_hash"`
	DisplayName         string    `bun:"display_name,notnull"`
	DiscordToken        string    `bun:"discord_token"`
	DiscordRefreshToken string    `bun:"discord_refresh_token"`
	DiscordUserID       string    `bun:"discord_user_id"`
	Profile             string    `bun:"profile"`
	AvatarURL           string    `bun:"avatar_url"`
	TwitterID           string    `bun:"twitter_id"`
	GithubID            string    `bun:"github_id"`
	CreatedAt           time.Time `bun:"created_at,notnull"`
	UpdatedAt           time.Time `bun:"updated_at,notnull"`
}
