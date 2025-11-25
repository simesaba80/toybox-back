//go:build integration

package token_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/testutil"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/token"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.Teardown()
	os.Exit(code)
}

func TestTokenRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := token.NewTokenRepository(db)

	ctx := context.Background()
	token := &entity.Token{
		RefreshToken: uuid.NewString(),
		UserID:       uuid.NewString(),
		ExpiredAt:    time.Now().Add(24 * time.Hour * 30),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	created, err := repo.Create(ctx, token)
	require.NoError(t, err)
	require.Equal(t, token.RefreshToken, created.RefreshToken)
	require.Equal(t, token.UserID, created.UserID)
	require.Equal(t, token.ExpiredAt, created.ExpiredAt)
	require.Equal(t, token.CreatedAt, created.CreatedAt)
	require.Equal(t, token.UpdatedAt, created.UpdatedAt)
}

func TestTokenRepository_CheckRefreshToken(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := token.NewTokenRepository(db)

	ctx := context.Background()
	token := &entity.Token{
		RefreshToken: uuid.NewString(),
		UserID:       uuid.NewString(),
		ExpiredAt:    time.Now().Add(24 * time.Hour * 30),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	created, err := repo.Create(ctx, token)
	require.NoError(t, err)

	userID, err := repo.CheckRefreshToken(ctx, created.RefreshToken)
	require.NoError(t, err)
	require.Equal(t, token.UserID, userID)

	userID, err = repo.CheckRefreshToken(ctx, "invalid")
	require.ErrorIs(t, err, domainerrors.ErrRefreshTokenInvalid)
	require.Empty(t, userID)

	expiredToken := &entity.Token{
		UserID: uuid.NewString(),
	}
	created, err = repo.Create(ctx, expiredToken)
	require.NoError(t, err)

	// 有効期限を過ぎたトークンを生成
	tokenDTO := new(dto.Token)
	db.NewUpdate().Model(tokenDTO).Where("refresh_token = ?", created.RefreshToken).Set("expired_at = ?", time.Now().Add(-24*time.Hour*60)).Exec(ctx)
	userID, err = repo.CheckRefreshToken(ctx, created.RefreshToken)
	require.ErrorIs(t, err, domainerrors.ErrRefreshTokenExpired)
	require.Empty(t, userID)
}
