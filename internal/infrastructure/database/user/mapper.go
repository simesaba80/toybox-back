package user

import (
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
)

func toDTO(e *entity.User) *dto.User {
	return &dto.User{
		ID:                  e.ID,
		Name:                e.Name,
		Email:               e.Email,
		PasswordHash:        e.PasswordHash,
		DisplayName:         e.DisplayName,
		DiscordToken:        e.DiscordToken,
		DiscordRefreshToken: e.DiscordRefreshToken,
		DiscordUserID:       e.DiscordUserID,
		Profile:             e.Profile,
		AvatarURL:           e.AvatarURL,
		TwitterID:           e.TwitterID,
		GithubID:            e.GithubID,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
	}
}

func toEntity(d *dto.User) *entity.User {
	return &entity.User{
		ID:                  d.ID,
		Name:                d.Name,
		Email:               d.Email,
		PasswordHash:        d.PasswordHash,
		DisplayName:         d.DisplayName,
		DiscordToken:        d.DiscordToken,
		DiscordRefreshToken: d.DiscordRefreshToken,
		DiscordUserID:       d.DiscordUserID,
		Profile:             d.Profile,
		AvatarURL:           d.AvatarURL,
		TwitterID:           d.TwitterID,
		GithubID:            d.GithubID,
		CreatedAt:           d.CreatedAt,
		UpdatedAt:           d.UpdatedAt,
	}
}

func toEntities(dtoUsers []*dto.User) []*entity.User {
	users := make([]*entity.User, len(dtoUsers))
	for i, d := range dtoUsers {
		users[i] = toEntity(d)
	}
	return users
}
