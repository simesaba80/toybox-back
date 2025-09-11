package entity

import (
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	ID        uuid.UUID
	WorkID    uuid.UUID
	AssetType AssetType
	UserID    uuid.UUID
	Extension string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
