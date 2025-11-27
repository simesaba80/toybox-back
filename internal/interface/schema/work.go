package schema

import (
	"time"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type GetWorkOutput struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	UserID      string          `json:"user_id"`
	Visibility  string          `json:"visibility"`
	Assets      []AssetResponse `json:"assets"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

type CreateWorkInput struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"required"`
	Visibility  string `json:"visibility"`
	UserID      string `json:"user_id" validate:"required"`
}

type CreateWorkOutput struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	Visibility  string `json:"visibility"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetWorksQuery struct {
	Limit *int `query:"limit" validate:"omitempty,min=1,max=100"`
	Page  *int `query:"page" validate:"omitempty,min=1"`
}

type WorkListResponse struct {
	Works      []GetWorkOutput `json:"works"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
}

type AssetResponse struct {
	ID        string `json:"id"`
	WorkID    string `json:"work_id"`
	AssetType string `json:"asset_type"`
	UserID    string `json:"user_id"`
	Extension string `json:"extension"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToWorkResponse(work *entity.Work) GetWorkOutput {
	if work == nil {
		return GetWorkOutput{}
	}
	return GetWorkOutput{
		ID:          work.ID.String(),
		Title:       work.Title,
		Description: work.Description,
		UserID:      work.UserID.String(),
		Visibility:  work.Visibility,
		Assets:      ToAssetResponses(work.Assets),
		CreatedAt:   work.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   work.UpdatedAt.Format(time.RFC3339),
	}
}

func ToCreateWorkOutput(work *entity.Work) CreateWorkOutput {
	if work == nil {
		return CreateWorkOutput{}
	}
	return CreateWorkOutput{
		ID:          work.ID.String(),
		Title:       work.Title,
		Description: work.Description,
		UserID:      work.UserID.String(),
		Visibility:  work.Visibility,
		CreatedAt:   work.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   work.UpdatedAt.Format(time.RFC3339),
	}
}

func ToAssetResponse(asset *entity.Asset) AssetResponse {
	if asset == nil {
		return AssetResponse{}
	}

	return AssetResponse{
		ID:        asset.ID,
		WorkID:    asset.WorkID,
		AssetType: asset.AssetType,
		UserID:    asset.UserID,
		Extension: asset.Extension,
		URL:       asset.URL,
		CreatedAt: asset.CreatedAt.Format(time.RFC3339),
		UpdatedAt: asset.UpdatedAt.Format(time.RFC3339),
	}
}

func ToAssetResponses(assets []*entity.Asset) []AssetResponse {
	if len(assets) == 0 {
		return []AssetResponse{}
	}

	res := make([]AssetResponse, 0, len(assets))
	for _, asset := range assets {
		res = append(res, ToAssetResponse(asset))
	}
	return res
}
