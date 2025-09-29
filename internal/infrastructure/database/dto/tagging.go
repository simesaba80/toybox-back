package dto

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Tagging struct {
	bun.BaseModel `bun:"table:tagging"`
	WorkID        uuid.UUID `json:"work_id" bun:"work_id,pk"`
	TagID         uuid.UUID `json:"tag_id" bun:"tag_id,pk"`
}
