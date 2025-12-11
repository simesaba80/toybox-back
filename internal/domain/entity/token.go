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

func NewToken(userID uuid.UUID) *Token {
	return &Token{
		RefreshToken: uuid.New(),
		UserID:       userID,
		ExpiredAt:    time.Now().Add(24 * time.Hour * 30),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
