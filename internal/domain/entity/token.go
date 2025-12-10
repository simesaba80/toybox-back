package entity

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	RefreshToken uuid.UUID
	UserID       uuid.UUID
	ExpiredAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
