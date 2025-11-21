package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
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

func (t *Token) ToTokenEntity() *entity.Token {
	return &entity.Token{
		RefreshToken: t.RefreshToken.String(),
		UserID:       t.UserID.String(),
		ExpiredAt:    t.ExpiredAt,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

func ToTokenDTO(entity *entity.Token) *Token {
	return &Token{
		RefreshToken: uuid.MustParse(entity.RefreshToken),
		UserID:       uuid.MustParse(entity.UserID),
		ExpiredAt:    entity.ExpiredAt,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}
