package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Token struct {
	bun.BaseModel `bun:"table:token"`
	RefreshToken  uuid.UUID `bun:"refresh_token,pk,default:gen_random_uuid()"`
	UserID        uuid.UUID `bun:"user_id,notnull"`
	ExpiredAt     time.Time `bun:"expired_at,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	UpdatedAt     time.Time `bun:"updated_at,notnull"`
}
