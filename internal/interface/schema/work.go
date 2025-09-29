package schema

import (
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type GetWorkOutput struct {
	ID              string         `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	UserID          string         `json:"user_id"`
	Visibility      string         `json:"visibility"`
	Assets          []entity.Asset `json:"assets"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}

type CreateWorkInput struct {
	Title           string `json:"title" validate:"required,max=100"`
	Description     string `json:"description" validate:"required"`
	Visibility      string `json:"visibility"`
	UserID          string `json:"user_id" validate:"required"`
}

type CreateWorkOutput struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	UserID          string `json:"user_id"`
	Visibility      string `json:"visibility"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type WorkListResponse struct {
	Works []GetWorkOutput `json:"works"`
}

func ToWorkResponse(work *entity.Work) GetWorkOutput {
	if work == nil {
		return GetWorkOutput{}
	}
	return GetWorkOutput{
		ID:              work.ID.String(),
		Title:           work.Title,
		Description:     work.Description,
		UserID:          work.UserID.String(),
		Visibility:      work.Visibility,
		Assets:          work.Assets,
		CreatedAt:       work.CreatedAt.String(),
		UpdatedAt:       work.UpdatedAt.String(),
	}
}

func ToCreateWorkOutput(work *entity.Work) CreateWorkOutput {
	if work == nil {
		return CreateWorkOutput{}
	}
	return CreateWorkOutput{
		ID:              work.ID.String(),
		Title:           work.Title,
		Description:     work.Description,
		UserID:          work.UserID.String(),
		Visibility:      work.Visibility,
		CreatedAt:       work.CreatedAt.String(),
		UpdatedAt:       work.UpdatedAt.String(),
	}
}
