package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Work struct {
	bun.BaseModel   `bun:"table:works"`
	ID              uuid.UUID `json:"id" bun:"id,pk"`
	Title           string    `json:"title" bun:"title,notnull" validate:"required,max=100"`
	Description     string    `json:"description" bun:"description,notnull" validate:"required"`
	DescriptionHTML string    `json:"description_html" bun:"description_html,notnull" validate:"required"`
	UserID          uuid.UUID `json:"user_id" bun:"user_id,notnull" validate:"required"`
	Visibility      string    `json:"visibility" bun:"visibility"`
	Assets          []Asset   `json:"assets" bun:"rel:has-many,join:id=work_id"`
	CreatedAt       time.Time `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt       time.Time `json:"updated_at" bun:"updated_at,notnull"`
}

func (w *Work) Validate() error {
	return validate.Struct(w)
}
