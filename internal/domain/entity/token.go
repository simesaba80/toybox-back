package entity

import "time"

type Token struct {
	RefreshToken string
	UserID       string
	ExpiredAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
