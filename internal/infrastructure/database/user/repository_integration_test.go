//go:build integration

package user_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/testutil"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/user"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.Teardown()
	os.Exit(code)
}

func TestUserRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := user.NewUserRepository(db)

	ctx := context.Background()
	user := &entity.User{
		ID:                  uuid.New(),
		Name:                "testuser",
		Email:               "testuser@example.com",
		DisplayName:         "testuser",
		AvatarURL:           "https://example.com/avatar.png",
		DiscordUserID:       "testuser",
		DiscordToken:        "testtoken",
		DiscordRefreshToken: "testrefreshtoken",
	}
	_, err := repo.Create(ctx, user)
	require.NoError(t, err)
}

func TestUserRepository_GetAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := user.NewUserRepository(db)

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		user := &entity.User{
			ID:            uuid.New(),
			Name:          "testuser" + strconv.Itoa(i),
			Email:         "testuser" + strconv.Itoa(i) + "@example.com",
			DisplayName:   "testuser",
			AvatarURL:     "https://example.com/avatar.png",
			DiscordUserID: "testuser",
		}
		_, err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	_, err := repo.GetAll(ctx)
	require.NoError(t, err)
}

func TestUserRepository_GetUserByDiscordUserID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := user.NewUserRepository(db)

	ctx := context.Background()
	user := &entity.User{
		Name:                "testuser",
		Email:               "testuser@example.com",
		DisplayName:         "testuser",
		AvatarURL:           "https://example.com/avatar.png",
		DiscordUserID:       "testuser",
		DiscordToken:        "testtoken",
		DiscordRefreshToken: "testrefreshtoken",
	}
	created, err := repo.Create(ctx, user)
	require.NoError(t, err)

	found, err := repo.GetUserByDiscordUserID(ctx, created.DiscordUserID)
	require.NoError(t, err)
	require.Equal(t, created.ID, found.ID)
	require.Equal(t, created.Name, found.Name)
	require.Equal(t, created.Email, found.Email)
	require.Equal(t, created.DisplayName, found.DisplayName)
	require.Equal(t, created.AvatarURL, found.AvatarURL)
	require.Equal(t, created.DiscordUserID, found.DiscordUserID)
	require.Equal(t, created.DiscordToken, found.DiscordToken)
	require.Equal(t, created.DiscordRefreshToken, found.DiscordRefreshToken)

	found, err = repo.GetUserByDiscordUserID(ctx, "notfound")
	require.ErrorIs(t, err, domainerrors.ErrUserNotFound)
	require.Nil(t, found)
}
