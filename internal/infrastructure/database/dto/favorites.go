package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Favorite struct {
	bun.BaseModel `bun:"table:favorites"`
	WorkID        uuid.UUID `json:"work_id" bun:"work_id,pk"`
	UserID        uuid.UUID `json:"user_id" bun:"user_id,pk"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,notnull"`
}
