package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/infrastructure/database/types"
)

type Work struct {
	bun.BaseModel   `bun:"table:work"`
	ID              uuid.UUID        `json:"id" bun:"id,pk"`
	Title           string           `json:"title" bun:"title,notnull"`
	Description     string           `json:"description" bun:"description,notnull"`
	UserID          uuid.UUID        `json:"user_id" bun:"user_id,notnull"`
	Visibility      types.Visibility `json:"visibility" bun:"visibility"`
	CreatedAt       time.Time        `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt       time.Time        `json:"updated_at" bun:"updated_at,notnull"`
}
