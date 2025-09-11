package dto

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Thumbnail struct {
	bun.BaseModel `bun:"table:thumbnails"`
	WorkID        uuid.UUID `json:"work_id" bun:"work_id,pk"`
	AssetID       uuid.UUID `json:"asset_id" bun:"asset_id,pk"`
}
