package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/types"
)

type Work struct {
	bun.BaseModel `bun:"table:work"`

	ID               uuid.UUID        `bun:"id,pk"`
	Title            string           `bun:"title,notnull"`
	Description      string           `bun:"description,notnull"`
	Visibility       types.Visibility `bun:"visibility"`
	ThumbnailAssetID uuid.UUID        `bun:"thumbnail_asset_id,notnull,rel:belongs-to,join:id=work_id"`
	Assets           []*Asset         `bun:"rel:has-many,join:id=work_id"`
	UserID           uuid.UUID        `bun:"user_id,notnull,rel:belongs-to,join:id=work_id"`
	CreatedAt        time.Time        `bun:"created_at,notnull"`
	UpdatedAt        time.Time        `bun:"updated_at,notnull"`
}

func (w *Work) ToWorkEntity() *entity.Work {
	assets := make([]*entity.Asset, len(w.Assets))
	for i, asset := range w.Assets {
		assets[i] = asset.ToAssetEntity()
	}

	return &entity.Work{
		ID:          w.ID.String(),
		Title:       w.Title,
		Description: w.Description,
		UserID:      w.UserID.String(),
		Visibility:  string(w.Visibility),
		Assets:      assets,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}

func ToWorkDTO(entity *entity.Work) *Work {
	assets := make([]*Asset, len(entity.Assets))
	for i, asset := range entity.Assets {
		assets[i] = ToAssetDTO(asset)
	}
	id := uuid.Nil
	userID := uuid.Nil
	thumbnailAssetID := uuid.Nil
	if entity.ID != "" {
		id = uuid.MustParse(entity.ID)
	}
	if entity.UserID != "" {
		userID = uuid.MustParse(entity.UserID)
	}
	if entity.ThumbnailAssetID != "" {
		thumbnailAssetID = uuid.MustParse(entity.ThumbnailAssetID)
	}

	return &Work{
		ID:               id,
		Title:            entity.Title,
		Description:      entity.Description,
		UserID:           userID,
		ThumbnailAssetID: thumbnailAssetID,
		Visibility:       types.Visibility(entity.Visibility),
		Assets:           assets,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}
