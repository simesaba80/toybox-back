package entity

import (
	"time"

	"github.com/google/uuid"
)

type Work struct {
	ID               uuid.UUID
	Title            string
	Description      string
	DescriptionHTML  string
	UserID           uuid.UUID
	Visibility       string
	ThumbnailAssetID uuid.UUID
	Assets           []*Asset
	URLs             []*string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewWork(title string, description string, userID uuid.UUID, visibility string, thumbnailAssetID uuid.UUID, assets []*Asset, urls []*string) *Work {
	return &Work{
		ID:               uuid.New(),
		Title:            title,
		Description:      description,
		UserID:           userID,
		Visibility:       visibility,
		ThumbnailAssetID: thumbnailAssetID,
		Assets:           assets,
		URLs:             urls,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
