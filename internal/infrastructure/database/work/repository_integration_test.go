//go:build integration

package work_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/testutil"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/work"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.Teardown()
	os.Exit(code)
}

func TestWorkRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := work.NewWorkRepository(db)

	ctx := context.Background()
	user := insertTestUser(t, db)
	tag1 := insertTestTag(t, db, "go")
	tag2 := insertTestTag(t, db, "rust")

	work := newTestWork(user.ID, "create-title")
	work.TagIDs = []uuid.UUID{tag1.ID, tag2.ID}
	work.Tags = []*entity.Tag{tag1, tag2}

	created, err := repo.Create(ctx, work)
	require.NoError(t, err)
	require.Equal(t, work.Title, created.Title)
	require.Equal(t, work.Description, created.Description)
	require.Equal(t, work.UserID, created.UserID)
	require.Equal(t, work.URLs, created.URLs)
	require.Equal(t, 2, len(created.Tags))
	require.Equal(t, 2, len(created.TagIDs))

	fetched, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
	require.Equal(t, 2, len(fetched.Tags))
}

func TestWorkRepository_GetAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := work.NewWorkRepository(db)

	ctx := context.Background()
	user := insertTestUser(t, db)
	tag := insertTestTag(t, db, "test")

	for i := 0; i < 3; i++ {
		work := newTestWork(user.ID, "title-"+uuid.NewString())
		work.TagIDs = []uuid.UUID{tag.ID}
		work.Tags = []*entity.Tag{tag}
		work.CreatedAt = work.CreatedAt.Add(time.Duration(i) * time.Minute)
		work.UpdatedAt = work.CreatedAt
		_, err := repo.Create(ctx, work)
		require.NoError(t, err)
	}

	works, total, err := repo.GetAll(ctx, 10, 0)
	require.NoError(t, err)
	require.Equal(t, 3, total)
	require.Len(t, works, 3)
	require.True(t, works[0].CreatedAt.After(works[1].CreatedAt) || works[0].CreatedAt.Equal(works[1].CreatedAt))

	for _, w := range works {
		require.Equal(t, 1, len(w.Tags))
		require.Equal(t, tag.Name, w.Tags[0].Name)
	}
}

func TestWorkRepository_GetByID_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := work.NewWorkRepository(db)

	_, err := repo.GetByID(context.Background(), uuid.New())
	require.ErrorIs(t, err, domainerrors.ErrWorkNotFound)
}

func TestWorkRepository_ExistsByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := work.NewWorkRepository(db)

	ctx := context.Background()
	user := insertTestUser(t, db)
	tag := insertTestTag(t, db, "test")

	work := newTestWork(user.ID, "exists-title")
	work.TagIDs = []uuid.UUID{tag.ID}
	work.Tags = []*entity.Tag{tag}
	created, err := repo.Create(ctx, work)
	require.NoError(t, err)

	exists, err := repo.ExistsById(ctx, created.ID)
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = repo.ExistsById(ctx, uuid.New())
	require.NoError(t, err)
	require.False(t, exists)
}

func insertTestUser(t *testing.T, db *bun.DB) *entity.User {
	t.Helper()

	now := time.Now().UTC().Truncate(time.Second)
	user := &entity.User{
		ID:            uuid.New(),
		Name:          "testuser",
		Email:         "testuser@example.com",
		DisplayName:   "testuser",
		DiscordUserID: "testuser",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	dtoUser := dto.ToUserDTO(user)
	_, err := db.NewInsert().Model(dtoUser).Exec(context.Background())
	require.NoError(t, err)

	return user
}

func newTestWork(userID uuid.UUID, title string) *entity.Work {
	now := time.Now().UTC().Truncate(time.Second)
	return &entity.Work{
		ID:               uuid.New(),
		Title:            title,
		Description:      "description",
		UserID:           userID,
		Visibility:       "public",
		ThumbnailAssetID: uuid.Nil,
		Assets:           []*entity.Asset{},
		URLs:             []*string{},
		TagIDs:           []uuid.UUID{},
		Tags:             []*entity.Tag{},
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func insertTestTag(t *testing.T, db *bun.DB, name string) *entity.Tag {
	t.Helper()

	now := time.Now().UTC().Truncate(time.Second)
	tag := &entity.Tag{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dtoTag := dto.ToTagDTO(tag)
	_, err := db.NewInsert().Model(dtoTag).Exec(context.Background())
	require.NoError(t, err)

	return tag
}

