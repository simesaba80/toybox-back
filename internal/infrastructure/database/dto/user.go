package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type User struct {
	bun.BaseModel `bun:"table:user"`

	ID                  uuid.UUID `bun:"id,pk,default:gen_random_uuid()"`
	Name                string    `bun:"name,notnull"`
	Email               string    `bun:"email,notnull,unique"`
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

func (u *User) ToUserEntity() *entity.User {
	return &entity.User{
		ID:                  u.ID,
		Name:                u.Name,
		Email:               u.Email,
		DisplayName:         u.DisplayName,
		DiscordToken:        u.DiscordToken,
		DiscordRefreshToken: u.DiscordRefreshToken,
		DiscordUserID:       u.DiscordUserID,
		Profile:             u.Profile,
		AvatarURL:           u.AvatarURL,
		TwitterID:           u.TwitterID,
		GithubID:            u.GithubID,
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
}

func ToUserDTO(entity *entity.User) *User {
	return &User{
		ID:                  entity.ID,
		Name:                entity.Name,
		Email:               entity.Email,
		DisplayName:         entity.DisplayName,
		DiscordToken:        entity.DiscordToken,
		DiscordRefreshToken: entity.DiscordRefreshToken,
		DiscordUserID:       entity.DiscordUserID,
		Profile:             entity.Profile,
		AvatarURL:           entity.AvatarURL,
		TwitterID:           entity.TwitterID,
		GithubID:            entity.GithubID,
		CreatedAt:           entity.CreatedAt,
		UpdatedAt:           entity.UpdatedAt,
	}
}
