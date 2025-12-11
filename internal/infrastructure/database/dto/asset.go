package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/types"
	"github.com/uptrace/bun"
)

type Asset struct {
	bun.BaseModel `bun:"table:asset"`
	ID            uuid.UUID       `bun:"id,pk"`
	WorkID        uuid.UUID       `bun:"work_id"`
	AssetType     types.AssetType `bun:"asset_type,notnull"`
	UserID        uuid.UUID       `bun:"user_id,notnull"`
	Extension     string          `bun:"extension,notnull"`
	URL           string          `bun:"url,notnull"`
	CreatedAt     time.Time       `bun:"created_at,notnull"`
	UpdatedAt     time.Time       `bun:"updated_at,notnull"`
}

func (a *Asset) ToAssetEntity() *entity.Asset {
	return &entity.Asset{
		ID:        a.ID,
		WorkID:    a.WorkID,
		UserID:    a.UserID,
		AssetType: string(a.AssetType),
		Extension: a.Extension,
		URL:       a.URL,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func ToAssetDTO(entity *entity.Asset) *Asset {

	return &Asset{
		ID:        entity.ID,
		WorkID:    entity.WorkID,
		UserID:    entity.UserID,
		AssetType: types.AssetType(entity.AssetType),
		Extension: entity.Extension,
		URL:       entity.URL,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
