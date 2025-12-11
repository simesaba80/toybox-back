package entity

import (
	"time"

	"github.com/google/uuid"
)

type Favorite struct {
	WorkID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}

func NewFavorite(workID uuid.UUID, userID uuid.UUID) *Favorite {
	return &Favorite{
		WorkID:    workID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
}
