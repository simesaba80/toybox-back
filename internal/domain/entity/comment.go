package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Comment struct {
	ID         uuid.UUID
	Content    string
	WorkID     uuid.UUID
	UserID     uuid.UUID
	ReplyAt    string
	User       *User
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (c *Comment) Validate() error {
	if c.Content == "" {
		return errors.New("content is required")
	}
	if utf8.RuneCountInString(c.Content) > 255 {
		return errors.New("content must be at most 255 characters")
	}
	if c.WorkID == uuid.Nil {
		return errors.New("work ID is required")
	}
	return nil
}