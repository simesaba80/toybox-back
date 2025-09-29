package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Tag struct {
	bun.BaseModel `bun:"table:tag"`
	ID            uuid.UUID `json:"id" bun:"id,pk"`
	Name          string    `json:"name" bun:"name,notnull"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,notnull"`
}
