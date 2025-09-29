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
	ID            uuid.UUID       `json:"id" bun:"id,pk"`
	WorkID        uuid.UUID       `json:"work_id" bun:"work_id"`
	AssetType     types.AssetType `json:"asset_type" bun:"asset_type,notnull"`
	UserID        uuid.UUID       `json:"user_id" bun:"user_id,notnull"`
	Extension     string          `json:"extension" bun:"extension,notnull"`
	URL           string          `json:"url" bun:"url,notnull"`
	CreatedAt     time.Time       `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt     time.Time       `json:"updated_at" bun:"updated_at,notnull"`
}

func (a *Asset) ToAssetEntity() *entity.Asset {
	return &entity.Asset{
		WorkID:    a.WorkID.String(),
		UserID:    a.UserID.String(),
		AssetType: string(a.AssetType),
		Extension: a.Extension,
		URL:       a.URL,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func ToAssetDTO(entity *entity.Asset) *Asset {
	workID, _ := uuid.Parse(entity.WorkID)
	userID, _ := uuid.Parse(entity.UserID)

	return &Asset{
		ID:        uuid.Nil,
		WorkID:    workID,
		UserID:    userID,
		AssetType: types.AssetType(entity.AssetType),
		Extension: entity.Extension,
		URL:       entity.URL,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
