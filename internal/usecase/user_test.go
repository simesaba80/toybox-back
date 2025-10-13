package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	tests := []struct {
		name         string
		inputName    string
		email        string
		passwordHash string
		displayName  string
		avatarURL    string
		setupMock    func(*usecase.MockUserRepository)
		wantErr      bool
	}{
		{
			name:         "正常系: ユーザー作成成功",
			inputName:    "testuser",
			email:        "test@example.com",
			passwordHash: "hashedpassword123",
			displayName:  "Test User",
			avatarURL:    "https://example.com/avatar.png",
			setupMock: func(m *usecase.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, user *entity.User) (*entity.User, error) {
						user.ID = uuid.New()
						user.CreatedAt = time.Now()
						user.UpdatedAt = time.Now()
						return user, nil
					}).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:         "異常系: バリデーションエラー(名前が空)",
			inputName:    "",
			email:        "test@example.com",
			passwordHash: "hashedpassword123",
			displayName:  "Test User",
			avatarURL:    "https://example.com/avatar.png",
			setupMock: func(m *usecase.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: true,
		},
		{
			name:         "異常系: バリデーションエラー(メールが空)",
			inputName:    "testuser",
			email:        "",
			passwordHash: "hashedpassword123",
			displayName:  "Test User",
			avatarURL:    "https://example.com/avatar.png",
			setupMock: func(m *usecase.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: true,
		},
		{
			name:         "異常系: バリデーションエラー(パスワードが空)",
			inputName:    "testuser",
			email:        "test@example.com",
			passwordHash: "",
			displayName:  "Test User",
			avatarURL:    "https://example.com/avatar.png",
			setupMock: func(m *usecase.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: true,
		},
		{
			name:         "異常系: バリデーションエラー(表示名が空)",
			inputName:    "testuser",
			email:        "test@example.com",
			passwordHash: "hashedpassword123",
			displayName:  "",
			avatarURL:    "https://example.com/avatar.png",
			setupMock: func(m *usecase.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: true,
		},
		{
			name:         "異常系: リポジトリエラー",
			inputName:    "testuser",
			email:        "test@example.com",
			passwordHash: "hashedpassword123",
			displayName:  "Test User",
			avatarURL:    "https://example.com/avatar.png",
			setupMock: func(m *usecase.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
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

			mockRepo := usecase.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)

			uc := usecase.NewUserUseCase(mockRepo, 30*time.Second)

			got, err := uc.CreateUser(context.Background(), tt.inputName, tt.email, tt.passwordHash, tt.displayName, tt.avatarURL)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.inputName, got.Name)
				assert.Equal(t, tt.email, got.Email)
				assert.Equal(t, tt.displayName, got.DisplayName)
				assert.False(t, got.CreatedAt.IsZero())
				assert.False(t, got.UpdatedAt.IsZero())
			}
		})
	}
}

func TestUserUseCase_GetAllUser(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*usecase.MockUserRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "正常系: ユーザー取得成功",
			setupMock: func(m *usecase.MockUserRepository) {
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
			setupMock: func(m *usecase.MockUserRepository) {
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
			setupMock: func(m *usecase.MockUserRepository) {
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

			mockRepo := usecase.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)

			uc := usecase.NewUserUseCase(mockRepo, 30*time.Second)

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
