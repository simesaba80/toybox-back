package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/uptrace/bun"
)

type Asset struct {
	bun.BaseModel `bun:"table:asset"`
	ID            uuid.UUID        `bun:"id,pk"`
	WorkID        uuid.UUID        `bun:"work_id"`
	AssetType     entity.AssetType `bun:"asset_type,notnull"`
	UserID        uuid.UUID        `bun:"user_id,notnull"`
	Extension     string           `bun:"extension,notnull"`
	URL           string           `bun:"url,notnull"`
	CreatedAt     time.Time        `bun:"created_at,notnull"`
	UpdatedAt     time.Time        `bun:"updated_at,notnull"`
}
