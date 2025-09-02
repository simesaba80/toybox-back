package dto

import (
	"time"

	"github.com/google/uuid"
)

type Work struct {
	ID              uuid.UUID `json:"id" bun:"id,pk"`
	Title           string    `json:"title" bun:"title,notnull"`
	Description     string    `json:"description" bun:"description,notnull"`
	DescriptionHTML string    `json:"description_html" bun:"description_html,notnull"`
	UserID          uuid.UUID `json:"user_id" bun:"user_id,notnull"`
	Visibility      string    `json:"visibility" bun:"visibility"`
	CreatedAt       time.Time `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt       time.Time `json:"updated_at" bun:"updated_at,notnull"`
}
