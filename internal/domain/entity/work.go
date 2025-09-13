package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Work struct {
	bun.BaseModel   `bun:"table:works"`
	ID              uuid.UUID `json:"id" bun:"id,pk"`
	Title           string    `json:"title" bun:"title,notnull"`
	Description     string    `json:"description" bun:"description,notnull"`
	DescriptionHTML string    `json:"description_html" bun:"description_html,notnull"`
	UserID          uuid.UUID `json:"user_id" bun:"user_id,notnull"`
	Visibility      string    `json:"visibility" bun:"visibility"`
	Assets          []Asset   `json:"assets" bun:"rel:has-many,join:id=work_id"`
	CreatedAt       time.Time `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt       time.Time `json:"updated_at" bun:"updated_at,notnull"`
}

func (w *Work) Validate() error {
	if w.Title == "" {
		return errors.New("title is required")
	}
	if utf8.RuneCountInString(w.Title) > 100 {
		return errors.New("title must be at most 100 characters")
	}
	if w.Description == "" {
		return errors.New("description is required")
	}
	if w.DescriptionHTML == "" {
		return errors.New("description_html is required")
	}
	if w.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	return nil
}