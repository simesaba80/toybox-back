package entity

import (
	"time"
)

type Asset struct {
	WorkID        string    `json:"work_id"`
	UserID        string    `json:"user_id"`
	AssetType     string    `json:"asset_type"`
	Extension     string    `json:"extension"`
	URL           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
