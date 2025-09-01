package entity

import (
	"errors"
	"time"
	"unicode/utf8"

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

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if utf8.RuneCountInString(u.Name) > 32 {
		return errors.New("name must be at most 32 characters")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if utf8.RuneCountInString(u.Email) > 256 {
		return errors.New("email must be at most 256 characters")
	}
	if u.PasswordHash == "" {
		return errors.New("password is required")
	}
	if u.DisplayName == "" {
		return errors.New("display name is required")
	}
	if utf8.RuneCountInString(u.DisplayName) > 32 {
		return errors.New("display name must be at most 32 characters")
	}
	if utf8.RuneCountInString(u.DiscordUserID) > 255 {
		return errors.New("discord user ID must be at most 255 characters")
	}
	if utf8.RuneCountInString(u.Profile) > 500 {
		return errors.New("profile must be at most 500 characters")
	}
	return nil
}
