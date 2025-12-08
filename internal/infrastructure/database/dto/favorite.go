package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/uptrace/bun"
)

type Favorite struct {
	bun.BaseModel `bun:"table:favorite"`
	WorkID        uuid.UUID `json:"work_id" bun:"work_id,pk"`
	UserID        uuid.UUID `json:"user_id" bun:"user_id,pk"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,notnull"`
}

func (f *Favorite) ToFavoriteEntity() *entity.Favorite {
	return &entity.Favorite{
		WorkID:    f.WorkID,
		UserID:    f.UserID,
		CreatedAt: f.CreatedAt,
	}
}

func ToFavoriteDTO(entity *entity.Favorite) *Favorite {
	return &Favorite{
		WorkID:    entity.WorkID,
		UserID:    entity.UserID,
		CreatedAt: entity.CreatedAt,
	}
}
