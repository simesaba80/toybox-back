package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Work struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	Visibility  string    `json:"visibility"`
	Assets      []Asset   `json:"assets"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	if w.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	return nil
}