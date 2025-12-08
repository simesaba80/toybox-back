package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/usecase"
	"github.com/simesaba80/toybox-back/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFavoriteUsecase_CreateFavorite(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mock.MockFavoriteRepository, uuid.UUID, uuid.UUID)
		wantErr   bool
		errIs     error
	}{
		{
			name: "正常系: 新規でいいねを作成できる",
			setupMock: func(m *mock.MockFavoriteRepository, workID, userID uuid.UUID) {
				m.EXPECT().
					Exists(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					DoAndReturn(func(_ context.Context, fav *entity.Favorite) bool {
						assert.Equal(t, workID, fav.WorkID)
						assert.Equal(t, userID, fav.UserID)
						return false
					})
				m.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					DoAndReturn(func(_ context.Context, fav *entity.Favorite) (*entity.Favorite, error) {
						assert.Equal(t, workID, fav.WorkID)
						assert.Equal(t, userID, fav.UserID)
						return fav, nil
					})
			},
			wantErr: false,
		},
		{
			name: "異常系: 既に存在するいいねは作成しない",
			setupMock: func(m *mock.MockFavoriteRepository, workID, userID uuid.UUID) {
				m.EXPECT().
					Exists(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					DoAndReturn(func(_ context.Context, fav *entity.Favorite) bool {
						assert.Equal(t, workID, fav.WorkID)
						assert.Equal(t, userID, fav.UserID)
						return true
					})
			},
			wantErr: true,
			errIs:   domainerrors.ErrFavoriteAlreadyExists,
		},
		{
			name: "異常系: リポジトリの作成エラーをラップして返す",
			setupMock: func(m *mock.MockFavoriteRepository, workID, userID uuid.UUID) {
				m.EXPECT().
					Exists(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					Return(false)
				m.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					Return(nil, domainerrors.ErrFailedToCreateFavorite)
			},
			wantErr: true,
			errIs:   domainerrors.ErrFailedToCreateFavorite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			workID := uuid.New()
			userID := uuid.New()

			mockRepo := mock.NewMockFavoriteRepository(ctrl)
			tt.setupMock(mockRepo, workID, userID)

			uc := usecase.NewFavoriteUsecase(mockRepo)

			err := uc.CreateFavorite(context.Background(), workID, userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errIs != nil {
					assert.ErrorIs(t, err, tt.errIs)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFavoriteUsecase_DeleteFavorite(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mock.MockFavoriteRepository, uuid.UUID, uuid.UUID)
		wantErr   bool
		errIs     error
	}{
		{
			name: "正常系: 既存のいいねを削除できる",
			setupMock: func(m *mock.MockFavoriteRepository, workID, userID uuid.UUID) {
				m.EXPECT().
					Exists(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					Return(true)
				m.EXPECT().
					Delete(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系: いいねが存在しない場合は削除しない",
			setupMock: func(m *mock.MockFavoriteRepository, workID, userID uuid.UUID) {
				m.EXPECT().
					Exists(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					Return(false)
			},
			wantErr: true,
			errIs:   domainerrors.ErrFavoriteNotFound,
		},
		{
			name: "異常系: 削除時のリポジトリエラーをそのまま返す",
			setupMock: func(m *mock.MockFavoriteRepository, workID, userID uuid.UUID) {
				m.EXPECT().
					Exists(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					Return(true)
				m.EXPECT().
					Delete(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
					Return(domainerrors.ErrFailedToDeleteFavorite)
			},
			wantErr: true,
			errIs:   domainerrors.ErrFailedToDeleteFavorite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			workID := uuid.New()
			userID := uuid.New()

			mockRepo := mock.NewMockFavoriteRepository(ctrl)
			tt.setupMock(mockRepo, workID, userID)

			uc := usecase.NewFavoriteUsecase(mockRepo)

			err := uc.DeleteFavorite(context.Background(), workID, userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errIs != nil {
					assert.ErrorIs(t, err, tt.errIs)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFavoriteUsecase_CountFavoritesByWorkID(t *testing.T) {
	workID := uuid.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFavoriteRepository(ctrl)

	gomock.InOrder(
		mockRepo.EXPECT().
			CountByWorkID(gomock.Any(), workID).
			Return(3, nil),
		mockRepo.EXPECT().
			CountByWorkID(gomock.Any(), workID).
			Return(0, domainerrors.ErrFailedToCountFavoritesByWorkID),
	)

	uc := usecase.NewFavoriteUsecase(mockRepo)

	total, err := uc.CountFavoritesByWorkID(context.Background(), workID)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)

	total, err = uc.CountFavoritesByWorkID(context.Background(), workID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domainerrors.ErrFailedToCountFavoritesByWorkID)
	assert.Equal(t, 0, total)
}

func TestFavoriteUsecase_IsFavorite(t *testing.T) {
	workID := uuid.New()
	userID := uuid.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockFavoriteRepository(ctrl)
	mockRepo.EXPECT().
		Exists(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
		Return(true)
	mockRepo.EXPECT().
		Exists(gomock.Any(), gomock.AssignableToTypeOf(&entity.Favorite{})).
		Return(false)

	uc := usecase.NewFavoriteUsecase(mockRepo)

	isFavorite := uc.IsFavorite(context.Background(), workID, userID)
	assert.True(t, isFavorite)

	isFavorite = uc.IsFavorite(context.Background(), workID, userID)
	assert.False(t, isFavorite)
}
