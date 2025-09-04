package entity

import "github.com/uptrace/bun"

type Asset struct {
	bun.BaseModel `bun:"table:asset"`
	WorkID        string `json:"work_id" bun:"work_id"`
	UserID        string `json:"user_id" bun:"user_id,notnull"`
	AssetType     string `json:"asset_type" bun:"asset_type,notnull"`
	Extension     string `json:"extension" bun:"extension,notnull"`
	URL           string `json:"url" bun:"url,notnull"`
	CreatedAt     string `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt     string `json:"updated_at" bun:"updated_at,notnull"`
}
