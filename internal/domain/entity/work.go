package entity

import (
	"time"
)

type Work struct {
	ID              string
	Title           string
	Description     string
	DescriptionHTML string
	UserID          string
	Visibility      string
	Assets          []*Asset
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
