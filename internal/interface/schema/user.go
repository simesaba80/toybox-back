package schema

import (
	"time"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type CreateUserInput struct {
	Name         string `json:"name" validate:"required,max=32"`
	Email        string `json:"email" validate:"required,email,max=256"`
	PasswordHash string `json:"password_hash" validate:"required"`
	DisplayName  string `json:"display_name" validate:"required,max=32"`
	AvatarURL    string `json:"avatar_url" validate:"omitempty,url"`
}

type CreateUserOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetUserOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Profile     string `json:"profile"`
	AvatarURL   string `json:"avatar_url"`
	TwitterID   string `json:"twitter_id"`
	GithubID    string `json:"github_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserListResponse struct {
	Users []GetUserOutput `json:"users"`
}

func ToUserResponse(user *entity.User) GetUserOutput {
	if user == nil {
		return GetUserOutput{}
	}
	return GetUserOutput{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Profile:     user.Profile,
		AvatarURL:   user.AvatarURL,
		TwitterID:   user.TwitterID,
		GithubID:    user.GithubID,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}
}
