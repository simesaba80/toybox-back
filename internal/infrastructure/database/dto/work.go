package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Work struct {
	bun.BaseModel `bun:"table:works"`

	ID              uuid.UUID `bun:"id,pk"`
	Title           string    `bun:"title,notnull"`
	Description     string    `bun:"description,notnull"`
	DescriptionHTML string    `bun:"description_html,notnull"`
	UserID          uuid.UUID `bun:"user_id,notnull"`
	Visibility      string    `bun:"visibility"`
	Assets          []*Asset  `bun:"rel:has-many,join:id=work_id"`
	CreatedAt       time.Time `bun:"created_at,notnull"`
	UpdatedAt       time.Time `bun:"updated_at,notnull"`
}
