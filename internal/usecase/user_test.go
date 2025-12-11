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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserUseCase_GetAllUser(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mock.MockUserRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "正常系: ユーザー取得成功",
			setupMock: func(m *mock.MockUserRepository) {
				expectedUsers := []*entity.User{
					{
						ID:          uuid.New(),
						Name:        "user1",
						Email:       "user1@example.com",
						DisplayName: "User One",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          uuid.New(),
						Name:        "user2",
						Email:       "user2@example.com",
						DisplayName: "User Two",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.EXPECT().
					GetAll(gomock.Any()).
					Return(expectedUsers, nil).
					Times(1)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "正常系: ユーザーが0件",
			setupMock: func(m *mock.MockUserRepository) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return([]*entity.User{}, nil).
					Times(1)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "異常系: リポジトリエラー",
			setupMock: func(m *mock.MockUserRepository) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return(nil, errors.New("database error")).
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

			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)

			uc := usecase.NewUserUseCase(mockRepo)

			got, err := uc.GetAllUser(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Len(t, got, tt.wantCount)
			}
		})
	}
}

func TestUserUseCase_GetByUserID(t *testing.T) {
	tests := []struct {
		name      string
		userID    uuid.UUID
		setupMock func(*mock.MockUserRepository, uuid.UUID)
		wantErr   bool
	}{
		{
			name:   "正常系: ユーザー取得成功",
			userID: uuid.New(),
			setupMock: func(m *mock.MockUserRepository, userID uuid.UUID) {
				expectedUser := &entity.User{
					ID:          userID,
					Name:        "testuser",
					Email:       "testuser@example.com",
					DisplayName: "testuser",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.EXPECT().
					GetByID(gomock.Any(), gomock.Eq(userID)).
					Return(expectedUser, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:   "異常系: リポジトリエラー",
			userID: uuid.New(),
			setupMock: func(m *mock.MockUserRepository, userID uuid.UUID) {
				m.EXPECT().
					GetByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo, tt.userID)

			uc := usecase.NewUserUseCase(mockRepo)

			got, err := uc.GetByUserID(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.userID, got.ID)
			}
		})
	}
}
