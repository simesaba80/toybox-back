package entity

import (
	"time"

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