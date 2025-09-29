package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Email               string    `json:"email"`
	PasswordHash        string    `json:"password_hash"`
	DisplayName         string    `json:"display_name"`
	DiscordToken        string    `json:"discord_token"`
	DiscordRefreshToken string    `json:"discord_refresh_token"`
	DiscordUserID       string    `json:"discord_user_id"`
	Profile             string    `json:"profile"`
	AvatarURL           string    `json:"avatar_url"`
	TwitterID           string    `json:"twitter_id"`
	GithubID            string    `json:"github_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
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