package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type Tag struct {
	bun.BaseModel `bun:"table:tag"`
	ID            uuid.UUID `json:"id" bun:"id,pk"`
	Name          string    `json:"name" bun:"name,notnull"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,notnull"`
}

func (t *Tag) ToTagEntity() *entity.Tag {
	return &entity.Tag{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func ToTagDTO(entity *entity.Tag) *Tag {
	return &Tag{
		ID:        entity.ID,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
