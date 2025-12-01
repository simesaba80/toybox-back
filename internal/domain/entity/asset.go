package entity

import (
	"time"
)

type Asset struct {
	ID        string
	WorkID    string
	AssetType string
	UserID    string
	Extension string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
