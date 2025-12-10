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
	dtoToken := dto.ToTokenDTO(token)

	_, err := r.db.NewInsert().Model(dtoToken).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return dtoToken.ToTokenEntity(), nil
}

func (r *TokenRepository) CheckRefreshToken(ctx context.Context, refreshToken uuid.UUID) (uuid.UUID, error) {
	dtoToken := new(dto.Token)
	err := r.db.NewSelect().Model(dtoToken).Where("refresh_token = ?", refreshToken).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, domainerrors.ErrRefreshTokenInvalid
		}
		return uuid.Nil, err
	}
	if dtoToken.ExpiredAt.Before(time.Now()) {
		return uuid.Nil, domainerrors.ErrRefreshTokenExpired
	}
	return dtoToken.UserID, nil
}

func (r *TokenRepository) UpdateRefreshToken(ctx context.Context, refreshToken uuid.UUID) (*entity.Token, error) {
	// 既存トークンを取得
	dtoToken := new(dto.Token)
	if err := r.db.NewSelect().Model(dtoToken).Where("refresh_token = ?", refreshToken).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerrors.ErrRefreshTokenInvalid
		}
		return nil, err
	}

	// 旧トークンを削除
	if _, err := r.db.NewDelete().Model(dtoToken).Where("refresh_token = ?", refreshToken).Exec(ctx); err != nil {
		return nil, err
	}

	// 新しいリフレッシュトークンを発行して再保存
	dtoToken.ExpiredAt = time.Now().Add(24 * time.Hour * 30)
	newRefreshToken := entity.NewToken(dtoToken.UserID)
	newRefreshToken, err := r.Create(ctx, newRefreshToken)
	if err != nil {
		return nil, err
	}
	return newRefreshToken, nil
}
