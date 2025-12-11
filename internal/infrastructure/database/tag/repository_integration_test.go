//go:build integration

package tag_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/tag"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/testutil"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.Teardown()
	os.Exit(code)
}

func TestTagRepository_ExistAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := tag.NewTagRepository(db)

	ctx := context.Background()
	tag1ID := insertTestTag(t, db, "go")
	tag2ID := insertTestTag(t, db, "rust")
	tag3ID := insertTestTag(t, db, "python")

	tests := []struct {
		name    string
		ids     []uuid.UUID
		want    bool
		wantErr bool
	}{
		{
			name: "正常系: 全てのタグが存在する",
			ids:  []uuid.UUID{tag1ID, tag2ID},
			want: true,
		},
		{
			name: "正常系: 1つのタグが存在する",
			ids:  []uuid.UUID{tag1ID},
			want: true,
		},
		{
			name: "正常系: 空のスライス",
			ids:  []uuid.UUID{},
			want: true,
		},
		{
			name: "異常系: 存在しないタグが含まれる",
			ids:  []uuid.UUID{tag1ID, uuid.New()},
			want: false,
		},
		{
			name: "異常系: 全て存在しない",
			ids:  []uuid.UUID{uuid.New(), uuid.New()},
			want: false,
		},
		{
			name: "正常系: 全てのタグが存在する（3つ）",
			ids:  []uuid.UUID{tag1ID, tag2ID, tag3ID},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistAll(ctx, tt.ids)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, exists)
			}
		})
	}
}

func TestTagRepository_FindAllByIDs(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := tag.NewTagRepository(db)

	ctx := context.Background()
	tag1ID := insertTestTag(t, db, "go")
	tag2ID := insertTestTag(t, db, "rust")
	tag3ID := insertTestTag(t, db, "python")

	tests := []struct {
		name      string
		ids       []uuid.UUID
		wantCount int
		wantNames []string
		wantErr   bool
	}{
		{
			name:      "正常系: 2つのタグを取得",
			ids:       []uuid.UUID{tag1ID, tag2ID},
			wantCount: 2,
			wantNames: []string{"go", "rust"},
		},
		{
			name:      "正常系: 1つのタグを取得",
			ids:       []uuid.UUID{tag3ID},
			wantCount: 1,
			wantNames: []string{"python"},
		},
		{
			name:      "正常系: 全てのタグを取得",
			ids:       []uuid.UUID{tag1ID, tag2ID, tag3ID},
			wantCount: 3,
			wantNames: []string{"go", "rust", "python"},
		},
		{
			name:      "正常系: 存在しないIDを含む（部分一致）",
			ids:       []uuid.UUID{tag1ID, uuid.New()},
			wantCount: 1,
			wantNames: []string{"go"},
		},
		{
			name:      "正常系: 全て存在しないID",
			ids:       []uuid.UUID{uuid.New(), uuid.New()},
			wantCount: 0,
			wantNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, err := repo.FindAllByIDs(ctx, tt.ids)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantCount, len(tags))

				// タグ名の確認
				tagNames := make([]string, len(tags))
				for i, tag := range tags {
					tagNames[i] = tag.Name
				}
				for _, wantName := range tt.wantNames {
					require.Contains(t, tagNames, wantName)
				}
			}
		})
	}
}

func insertTestTag(t *testing.T, db *bun.DB, name string) uuid.UUID {
	t.Helper()

	now := time.Now().UTC().Truncate(time.Second)
	tag := &dto.Tag{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := db.NewInsert().Model(tag).Exec(context.Background())
	require.NoError(t, err)

	return tag.ID
}
