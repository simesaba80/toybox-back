package entity

import (
	"time"

	"github.com/google/uuid"
)

type Work struct {
	ID              uuid.UUID
	Title           string
	Description     string
	DescriptionHTML string
	UserID          uuid.UUID
	Visibility      string
	Assets          []*Asset
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
