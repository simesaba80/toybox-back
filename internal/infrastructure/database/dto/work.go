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

	ID          uuid.UUID `bun:"id,pk"`
	Title       string    `bun:"title,notnull"`
	Description string    `bun:"description,notnull"`
	UserID      uuid.UUID `bun:"user_id,notnull"`
	Visibility  string    `bun:"visibility"`
	Assets      []*Asset  `bun:"rel:has-many,join:id=work_id"`
	CreatedAt   time.Time `bun:"created_at,notnull"`
	UpdatedAt   time.Time `bun:"updated_at,notnull"`
}

func (w *Work) ToWorkEntity() *entity.Work {
	return &entity.Work{
		ID:          w.ID,
		Title:       w.Title,
		Description: w.Description,
		UserID:      w.UserID,
		Visibility:  string(w.Visibility),
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}

func ToWorkDTO(entity *entity.Work) *Work {
	return &Work{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		UserID:      entity.UserID,
		Visibility:  types.Visibility(entity.Visibility),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
