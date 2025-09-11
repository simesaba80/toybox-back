package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/infrastructure/database/types"
)

type Comment struct {
	bun.BaseModel `bun:"table:comment"`
	ID            uuid.UUID        `json:"id" bun:"id,pk"`
	Content       string           `json:"content" bun:"content,notnull"`
	WorkID        uuid.UUID        `json:"work_id" bun:"work_id,notnull"`
	UserID        uuid.UUID        `json:"user_id" bun:"user_id"`
	ReplyAt       *uuid.UUID       `json:"reply_at" bun:"reply_at"`
	Visibility    types.Visibility `json:"visibility" bun:"visibility,notnull"`
	CreatedAt     time.Time        `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt     time.Time        `json:"updated_at" bun:"updated_at,notnull"`
}
