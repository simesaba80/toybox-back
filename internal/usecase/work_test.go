package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/usecase"
	"github.com/simesaba80/toybox-back/internal/usecase/mock"
	"github.com/simesaba80/toybox-back/internal/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorkUseCase_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		limit         *int
		page          *int
		userID        uuid.UUID
		setupWorkMock func(*mock.MockWorkRepository)
		setupTagMock  func(*mock.MockTagRepository)
		wantCount     int
		wantTotal     int
		wantLimit     int
		wantPage      int
		wantErr       bool
	}{
		{
			name:   "正常系: デフォルトページネーション",
			limit:  nil,
			page:   nil,
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				expectedWorks := []*entity.Work{
					{ID: uuid.New(), Title: "Work1", Description: "Desc1", UserID: uuid.New()},
					{ID: uuid.New(), Title: "Work2", Description: "Desc2", UserID: uuid.New()},
				}
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(20), gomock.Eq(0)).
					Return(expectedWorks, 50, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    2,
			wantTotal:    50,
			wantLimit:    20,
			wantPage:     1,
			wantErr:      false,
		},
		{
			name:   "正常系: カスタムページネーション(limit=10, page=1)",
			limit:  util.IntPtr(10),
			page:   util.IntPtr(1),
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				expectedWorks := []*entity.Work{
					{ID: uuid.New(), Title: "Work1", Description: "Desc1", UserID: uuid.New()},
				}
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(10), gomock.Eq(0)).
					Return(expectedWorks, 30, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    1,
			wantTotal:    30,
			wantLimit:    10,
			wantPage:     1,
			wantErr:      false,
		},
		{
			name:   "正常系: カスタムページネーション(limit=20, page=2)",
			limit:  util.IntPtr(20),
			page:   util.IntPtr(2),
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				expectedWorks := []*entity.Work{
					{ID: uuid.New(), Title: "Work3", Description: "Desc3", UserID: uuid.New()},
				}
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(20), gomock.Eq(20)).
					Return(expectedWorks, 50, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    1,
			wantTotal:    50,
			wantLimit:    20,
			wantPage:     2,
			wantErr:      false,
		},
		{
			name:   "正常系: 作品が0件",
			limit:  nil,
			page:   nil,
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(20), gomock.Eq(0)).
					Return([]*entity.Work{}, 0, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    0,
			wantTotal:    0,
			wantLimit:    20,
			wantPage:     1,
			wantErr:      false,
		},
		{
			name:   "エッジケース: limit=0, page=0",
			limit:  util.IntPtr(0),
			page:   util.IntPtr(0),
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(0), gomock.Eq(0)).
					Return([]*entity.Work{}, 0, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    0,
			wantTotal:    0,
			wantLimit:    0,
			wantPage:     0,
			wantErr:      false,
		},
		{
			name:   "エッジケース: 負の値(limit=-1, page=-1)",
			limit:  util.IntPtr(-1),
			page:   util.IntPtr(-1),
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(-1), gomock.Eq(2)).
					Return([]*entity.Work{}, 0, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    0,
			wantTotal:    0,
			wantLimit:    -1,
			wantPage:     -1,
			wantErr:      false,
		},
		{
			name:   "エッジケース: limitのみ指定、pageはnil",
			limit:  util.IntPtr(5),
			page:   nil,
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				expectedWorks := []*entity.Work{
					{ID: uuid.New(), Title: "Work1", Description: "Desc1", UserID: uuid.New()},
				}
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(5), gomock.Eq(0)).
					Return(expectedWorks, 10, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    1,
			wantTotal:    10,
			wantLimit:    5,
			wantPage:     1,
			wantErr:      false,
		},
		{
			name:   "エッジケース: pageのみ指定、limitはnil",
			limit:  nil,
			page:   util.IntPtr(3),
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				expectedWorks := []*entity.Work{
					{ID: uuid.New(), Title: "Work1", Description: "Desc1", UserID: uuid.New()},
				}
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(20), gomock.Eq(40)).
					Return(expectedWorks, 100, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    1,
			wantTotal:    100,
			wantLimit:    20,
			wantPage:     3,
			wantErr:      false,
		},
		{
			name:   "異常系: リポジトリエラー",
			limit:  nil,
			page:   nil,
			userID: uuid.Nil,
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					GetAllPublic(gomock.Any(), gomock.Eq(20), gomock.Eq(0)).
					Return(nil, 0, errors.New("database connection failed")).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    0,
			wantTotal:    0,
			wantLimit:    0,
			wantPage:     0,
			wantErr:      true,
		},
		{
			name:   "正常系: 認証済みユーザーは限定作品含め取得",
			limit:  nil,
			page:   util.IntPtr(2),
			userID: uuid.New(),
			setupWorkMock: func(m *mock.MockWorkRepository) {
				expectedWorks := []*entity.Work{
					{ID: uuid.New(), Title: "PrivateWork", Description: "Desc", UserID: uuid.New()},
				}
				m.EXPECT().
					GetAll(gomock.Any(), gomock.Eq(20), gomock.Eq(20)).
					Return(expectedWorks, 30, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantCount:    1,
			wantTotal:    30,
			wantLimit:    20,
			wantPage:     2,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWorkRepo := mock.NewMockWorkRepository(ctrl)
			mockTagRepo := mock.NewMockTagRepository(ctrl)
			tt.setupWorkMock(mockWorkRepo)
			tt.setupTagMock(mockTagRepo)

			uc := usecase.NewWorkUseCase(mockWorkRepo, mockTagRepo)

			got, total, limit, page, err := uc.GetAll(context.Background(), tt.limit, tt.page, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Len(t, got, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
				assert.Equal(t, tt.wantLimit, limit)
				assert.Equal(t, tt.wantPage, page)
			}
		})
	}
}

func TestWorkUseCase_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		workID        uuid.UUID
		setupWorkMock func(*mock.MockWorkRepository, uuid.UUID)
		setupTagMock  func(*mock.MockTagRepository)
		wantErr       bool
	}{
		{
			name:   "正常系: 作品取得成功",
			workID: uuid.New(),
			setupWorkMock: func(m *mock.MockWorkRepository, workID uuid.UUID) {
				expectedWork := &entity.Work{
					ID:          workID,
					Title:       "Test Work",
					Description: "Test Description",
					UserID:      uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.EXPECT().
					GetByID(gomock.Any(), gomock.Eq(workID)).
					Return(expectedWork, nil).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantErr:      false,
		},
		{
			name:   "異常系: リポジトリエラー",
			workID: uuid.New(),
			setupWorkMock: func(m *mock.MockWorkRepository, workID uuid.UUID) {
				m.EXPECT().
					GetByID(gomock.Any(), gomock.Eq(workID)).
					Return(nil, errors.New("work not found")).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository) {},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWorkRepo := mock.NewMockWorkRepository(ctrl)
			mockTagRepo := mock.NewMockTagRepository(ctrl)
			tt.setupWorkMock(mockWorkRepo, tt.workID)
			tt.setupTagMock(mockTagRepo)

			uc := usecase.NewWorkUseCase(mockWorkRepo, mockTagRepo)

			got, err := uc.GetByID(context.Background(), tt.workID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.workID, got.ID)
			}
		})
	}
}

func TestWorkUseCase_GetByUserID(t *testing.T) {
	targetUserID := uuid.New()
	authenticatedUserID := uuid.New()

	tests := []struct {
		name                string
		userID              uuid.UUID
		authenticatedUserID uuid.UUID
		setupMock           func(*mock.MockWorkRepository, uuid.UUID)
		wantCount           int
		wantErr             bool
	}{
		{
			name:                "正常系: 認証済みユーザー（公開・非公開両方取得）",
			userID:              targetUserID,
			authenticatedUserID: authenticatedUserID,
			setupMock: func(m *mock.MockWorkRepository, userID uuid.UUID) {
				expectedWorks := []*entity.Work{
					{
						ID:          uuid.New(),
						Title:       "Public Work",
						Description: "Public Description",
						UserID:      userID,
						Visibility:  "public",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          uuid.New(),
						Title:       "Private Work",
						Description: "Private Description",
						UserID:      userID,
						Visibility:  "private",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.EXPECT().
					GetByUserID(gomock.Any(), gomock.Eq(userID), gomock.Eq(false)).
					Return(expectedWorks, nil).
					Times(1)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:                "正常系: 未認証ユーザー（公開作品のみ取得）",
			userID:              targetUserID,
			authenticatedUserID: uuid.Nil,
			setupMock: func(m *mock.MockWorkRepository, userID uuid.UUID) {
				expectedWorks := []*entity.Work{
					{
						ID:          uuid.New(),
						Title:       "Public Work",
						Description: "Public Description",
						UserID:      userID,
						Visibility:  "public",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.EXPECT().
					GetByUserID(gomock.Any(), gomock.Eq(userID), gomock.Eq(true)).
					Return(expectedWorks, nil).
					Times(1)
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:                "正常系: 作品が0件",
			userID:              targetUserID,
			authenticatedUserID: uuid.Nil,
			setupMock: func(m *mock.MockWorkRepository, userID uuid.UUID) {
				m.EXPECT().
					GetByUserID(gomock.Any(), gomock.Eq(userID), gomock.Eq(true)).
					Return([]*entity.Work{}, nil).
					Times(1)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:                "異常系: リポジトリエラー（認証済み）",
			userID:              targetUserID,
			authenticatedUserID: authenticatedUserID,
			setupMock: func(m *mock.MockWorkRepository, userID uuid.UUID) {
				m.EXPECT().
					GetByUserID(gomock.Any(), gomock.Eq(userID), gomock.Eq(false)).
					Return(nil, errors.New("database connection failed")).
					Times(1)
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:                "異常系: リポジトリエラー（未認証）",
			userID:              targetUserID,
			authenticatedUserID: uuid.Nil,
			setupMock: func(m *mock.MockWorkRepository, userID uuid.UUID) {
				m.EXPECT().
					GetByUserID(gomock.Any(), gomock.Eq(userID), gomock.Eq(true)).
					Return(nil, errors.New("database connection failed")).
					Times(1)
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockWorkRepository(ctrl)
			mockTagRepo := mock.NewMockTagRepository(ctrl)
			tt.setupMock(mockRepo, tt.userID)
			tt.setupTagMock(mockTagRepo)

			uc := usecase.NewWorkUseCase(mockRepo, mockTagRepo)

			got, err := uc.GetByUserID(context.Background(), tt.userID, tt.authenticatedUserID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Len(t, got, tt.wantCount)
				if tt.wantCount > 0 {
					for _, work := range got {
						assert.Equal(t, tt.userID, work.UserID)
					}
				}
			}
		})
	}
}

func TestWorkUseCase_CreateWork(t *testing.T) {
	tests := []struct {
		name             string
		title            string
		description      string
		visibility       string
		thumbnailAssetID uuid.UUID
		assetIDs         []uuid.UUID
		urls             []string
		userID           uuid.UUID
		tagIDs           []uuid.UUID
		setupWorkMock    func(*mock.MockWorkRepository)
		setupTagMock     func(*mock.MockTagRepository, []uuid.UUID)
		wantErr          bool
	}{
		{
			name:             "正常系: 作品作成成功",
			title:            "New Work",
			description:      "New Description",
			visibility:       "public",
			thumbnailAssetID: uuid.New(),
			assetIDs:         []uuid.UUID{uuid.New()},
			urls:             []string{"https://example.com"},
			userID:           uuid.New(),
			tagIDs:           []uuid.UUID{uuid.New(), uuid.New()},
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, work *entity.Work) (*entity.Work, error) {
						work.CreatedAt = time.Now()
						work.UpdatedAt = time.Now()
						return work, nil
					}).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository, tagIDs []uuid.UUID) {
				m.EXPECT().
					ExistAll(gomock.Any(), gomock.Eq(tagIDs)).
					Return(true, nil).
					Times(1)
				m.EXPECT().
					FindAllByIDs(gomock.Any(), gomock.Eq(tagIDs)).
					Return([]*entity.Tag{
						{ID: tagIDs[0], Name: "Tag1"},
						{ID: tagIDs[1], Name: "Tag2"},
					}, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:             "異常系: バリデーションエラー(タイトル空)",
			title:            "",
			description:      "Description",
			visibility:       "public",
			thumbnailAssetID: uuid.New(),
			assetIDs:         []uuid.UUID{uuid.New()},
			urls:             []string{"https://example.com"},
			userID:           uuid.New(),
			tagIDs:           []uuid.UUID{},
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupTagMock: func(m *mock.MockTagRepository, tagIDs []uuid.UUID) {},
			wantErr:      true,
		},
		{
			name:             "異常系: バリデーションエラー(説明空)",
			title:            "Title",
			description:      "",
			visibility:       "public",
			thumbnailAssetID: uuid.New(),
			assetIDs:         []uuid.UUID{uuid.New()},
			urls:             []string{"https://example.com"},
			userID:           uuid.New(),
			tagIDs:           []uuid.UUID{},
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupTagMock: func(m *mock.MockTagRepository, tagIDs []uuid.UUID) {},
			wantErr:      true,
		},
		{
			name:             "異常系: バリデーションエラー(可視性空)",
			title:            "Title",
			description:      "Description",
			visibility:       "",
			thumbnailAssetID: uuid.New(),
			assetIDs:         []uuid.UUID{uuid.New()},
			urls:             []string{"https://example.com"},
			userID:           uuid.New(),
			tagIDs:           []uuid.UUID{uuid.New()},
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupTagMock: func(m *mock.MockTagRepository, tagIDs []uuid.UUID) {},
			wantErr:      true,
		},
		{
			name:             "異常系: バリデーションエラー(タグなし)",
			title:            "Title",
			description:      "Description",
			visibility:       "public",
			thumbnailAssetID: uuid.New(),
			assetIDs:         []uuid.UUID{uuid.New()},
			urls:             []string{"https://example.com"},
			userID:           uuid.New(),
			tagIDs:           []uuid.UUID{},
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupTagMock: func(m *mock.MockTagRepository, tagIDs []uuid.UUID) {},
			wantErr:      true,
		},
		{
			name:             "異常系: タグが存在しない",
			title:            "New Work",
			description:      "New Description",
			visibility:       "public",
			thumbnailAssetID: uuid.New(),
			assetIDs:         []uuid.UUID{uuid.New()},
			urls:             []string{"https://example.com"},
			userID:           uuid.New(),
			tagIDs:           []uuid.UUID{uuid.New()},
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupTagMock: func(m *mock.MockTagRepository, tagIDs []uuid.UUID) {
				m.EXPECT().
					ExistAll(gomock.Any(), gomock.Eq(tagIDs)).
					Return(false, nil).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:             "異常系: リポジトリエラー",
			title:            "New Work",
			description:      "New Description",
			visibility:       "public",
			thumbnailAssetID: uuid.New(),
			assetIDs:         []uuid.UUID{uuid.New()},
			urls:             []string{"https://example.com"},
			userID:           uuid.New(),
			tagIDs:           []uuid.UUID{uuid.New()},
			setupWorkMock: func(m *mock.MockWorkRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			setupTagMock: func(m *mock.MockTagRepository, tagIDs []uuid.UUID) {
				m.EXPECT().
					ExistAll(gomock.Any(), gomock.Eq(tagIDs)).
					Return(true, nil).
					Times(1)
				m.EXPECT().
					FindAllByIDs(gomock.Any(), gomock.Eq(tagIDs)).
					Return([]*entity.Tag{{ID: tagIDs[0], Name: "Tag1"}}, nil).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWorkRepo := mock.NewMockWorkRepository(ctrl)
			mockTagRepo := mock.NewMockTagRepository(ctrl)
			tt.setupWorkMock(mockWorkRepo)
			tt.setupTagMock(mockTagRepo, tt.tagIDs)

			uc := usecase.NewWorkUseCase(mockWorkRepo, mockTagRepo)
			got, err := uc.CreateWork(context.Background(), tt.title, tt.description, tt.visibility, tt.thumbnailAssetID, tt.assetIDs, tt.urls, tt.userID, tt.tagIDs)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.title, got.Title)
				assert.Equal(t, tt.description, got.Description)
				assert.Equal(t, tt.userID, got.UserID)
			}
		})
	}
}
