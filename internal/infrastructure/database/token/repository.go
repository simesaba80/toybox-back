package token

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

type TokenRepository struct {
	db *bun.DB
}

func NewTokenRepository(db *bun.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) Create(ctx context.Context, token *entity.Token) (*entity.Token, error) {
	now := time.Now()

	token.RefreshToken = uuid.NewString()
	token.ExpiredAt = now.Add(24 * time.Hour * 30)
	token.CreatedAt = now
	token.UpdatedAt = now

	dtoToken := dto.ToTokenDTO(token)

	_, err := r.db.NewInsert().Model(dtoToken).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return dtoToken.ToTokenEntity(), nil
}

func (r *TokenRepository) CheckRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	dtoToken := new(dto.Token)
	err := r.db.NewSelect().Model(dtoToken).Where("refresh_token = ?", refreshToken).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domainerrors.ErrRefreshTokenInvalid
		}
		return "", err
	}
	if dtoToken.ExpiredAt.Before(time.Now()) {
		return "", domainerrors.ErrRefreshTokenExpired
	}
	return dtoToken.UserID.String(), nil
}
