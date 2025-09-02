package schema

import (
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type GetWorkOutput struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DescriptionHTML string `json:"description_html"`
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
		DescriptionHTML: work.DescriptionHTML,
		UserID:          work.UserID.String(),
		Visibility:      work.Visibility,
		CreatedAt:       work.CreatedAt.String(),
		UpdatedAt:       work.UpdatedAt.String(),
	}
}
