package entity

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user"`
	ID                  uuid.UUID `json:"id" bun:"id,pk"`
	Name                string    `json:"name" bun:"name,notnull" validate:"required,max=32"`
	Email               string    `json:"email" bun:"email,notnull,unique" validate:"required,email,max=256"`
	PasswordHash        string    `json:"password_hash" bun:"password_hash" validate:"required"`
	DisplayName         string    `json:"display_name" bun:"display_name,notnull" validate:"required,max=32"`
	DiscordToken        string    `json:"discord_token" bun:"discord_token"`
	DiscordRefreshToken string    `json:"discord_refresh_token" bun:"discord_refresh_token"`
	DiscordUserID       string    `json:"discord_user_id" bun:"discord_user_id" validate:"max=255"`
	Profile             string    `json:"profile" bun:"profile" validate:"max=500"`
	AvatarURL           string    `json:"avatar_url" bun:"avatar_url" validate:"omitempty,url"`
	TwitterID           string    `json:"twitter_id" bun:"twitter_id"`
	GithubID            string    `json:"github_id" bun:"github_id"`
	CreatedAt           time.Time `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt           time.Time `json:"updated_at" bun:"updated_at,notnull"`
}

var validate = validator.New()

func (u *User) Validate() error {
	return validate.Struct(u)
}
