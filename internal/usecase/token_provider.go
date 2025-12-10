package usecase

import "github.com/google/uuid"

type TokenProvider interface {
	GenerateToken(userID uuid.UUID) (string, error)
}
