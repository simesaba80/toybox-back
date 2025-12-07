package entity

import "github.com/google/uuid"

type Favorite struct {
	WorkID uuid.UUID
	UserID uuid.UUID
}

func NewFavorite(workID uuid.UUID, userID uuid.UUID) *Favorite {
	return &Favorite{
		WorkID: workID,
		UserID: userID,
	}
}
