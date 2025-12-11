package entity

import (
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	ID        uuid.UUID
	WorkID    uuid.UUID
	AssetType string
	UserID    uuid.UUID
	Extension string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAsset(assetType string, userID uuid.UUID, extension string, url string) *Asset {
	return &Asset{
		ID:        uuid.New(),
		WorkID:    uuid.Nil,
		AssetType: assetType,
		UserID:    userID,
		Extension: extension,
		URL:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
