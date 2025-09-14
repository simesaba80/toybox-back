package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/types"
	"github.com/uptrace/bun"
)

type URLInfo struct {
	bun.BaseModel `bun:"table:url_info"`
	ID            uuid.UUID     `json:"id" bun:"id,pk"`
	WorkID        uuid.UUID     `json:"work_id" bun:"work_id,notnull"`
	URL           string        `json:"url" bun:"url,notnull"`
	URLType       types.URLType `json:"url_type" bun:"url_type,notnull"`
	UserID        uuid.UUID     `json:"user_id" bun:"user_id,notnull"`
	CreatedAt     time.Time     `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt     time.Time     `json:"updated_at" bun:"updated_at,notnull"`
}
