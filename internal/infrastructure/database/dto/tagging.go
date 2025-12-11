package dto

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Tagging struct {
	bun.BaseModel `bun:"table:tagging"`
	WorkID        uuid.UUID `json:"work_id" bun:"work_id,pk"`
	TagID         uuid.UUID `json:"tag_id" bun:"tag_id,pk"`
	Work          *Work     `bun:"rel:belongs-to,join:work_id=id"`
	Tag           *Tag      `bun:"rel:belongs-to,join:tag_id=id"`
}
