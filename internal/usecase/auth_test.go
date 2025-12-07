package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthUsecase_GetDiscordAuthURL(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*mock.MockDiscordRepository, *mock.MockUserRepository, *mock.MockTokenProvider, *mock.MockTokenRepository)
		wantAuthURL string
		wantErr     bool
	}{
		{
			name: "正常系: ディスコード認証URL取得成功",
			setupMock: func(m *mock.MockDiscordRepository, _ *mock.MockUserRepository, _ *mock.MockTokenProvider, _ *mock.MockTokenRepository) {
				m.EXPECT().
					GetDiscordClientID(gomock.Any()).
					Return("1234567890", nil).
					Times(1)
				m.EXPECT().
					GetHostURL(gomock.Any()).
					Return("https://localhost:8080", nil).
					Times(1)
				m.EXPECT().
					GetDiscordAuthURL(gomock.Any()).
					Return("https://localhost:8080/auth/discord/callback", nil).
					Times(1)
			},
			wantAuthURL: "https://localhost:8080/auth/discord/callback",
			wantErr:     false,
		},
		{
			name: "異常系: ディスコードクライアントID取得失敗",
			setupMock: func(m *mock.MockDiscordRepository, _ *mock.MockUserRepository, _ *mock.MockTokenProvider, _ *mock.MockTokenRepository) {
				m.EXPECT().
					GetDiscordClientID(gomock.Any()).
					Return("", domainerrors.ErrClientIDNotSet).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "異常系: リダイレクトURL取得失敗",
			setupMock: func(m *mock.MockDiscordRepository, _ *mock.MockUserRepository, _ *mock.MockTokenProvider, _ *mock.MockTokenRepository) {
				m.EXPECT().
					GetDiscordClientID(gomock.Any()).
					Return("1234567890", nil).
					Times(1)
				m.EXPECT().
					GetHostURL(gomock.Any()).
					Return("", domainerrors.ErrRedirectURLNotSet).
					Times(1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDiscordRepository := mock.NewMockDiscordRepository(ctrl)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockTokenProvider := mock.NewMockTokenProvider(ctrl)
			mockTokenRepository := mock.NewMockTokenRepository(ctrl)
			tt.setupMock(mockDiscordRepository, mockUserRepository, mockTokenProvider, mockTokenRepository)
			uc := NewAuthUsecase(mockDiscordRepository, mockUserRepository, mockTokenProvider, mockTokenRepository)
			got, err := uc.GetDiscordAuthURL(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			}
		})
	}
}

func TestAuthUsecase_AuthenticateUser(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mock.MockDiscordRepository, *mock.MockUserRepository, *mock.MockTokenProvider, *mock.MockTokenRepository)
		wantErr   bool
	}{
		{
			name: "正常系: 既存ユーザーのディスコード認証成功",
			setupMock: func(m *mock.MockDiscordRepository, u *mock.MockUserRepository, tp *mock.MockTokenProvider, tr *mock.MockTokenRepository) {
				m.EXPECT().
					GetDiscordToken(gomock.Any(), gomock.Any()).
					Return(entity.DiscordToken{
						AccessToken:  "test",
						RefreshToken: "test",
						Expiry:       time.Now().Add(1 * time.Hour),
						ExpiresIn:    3600,
						TokenType:    "Bearer",
					}, nil).
					Times(1)
				m.EXPECT().
					FetchDiscordUser(gomock.Any(), gomock.Any()).
					Return(entity.DiscordUser{
						ID: "test",
					}, nil).
					Times(1)
				m.EXPECT().
					GetDiscordGuilds(gomock.Any(), gomock.Any()).
					Return([]string{"test"}, nil).
					Times(1)
				m.EXPECT().
					GetAllowedDiscordGuilds(gomock.Any()).
					Return([]string{"test"}, nil).
					Times(1)
				u.EXPECT().
					GetUserByDiscordUserID(gomock.Any(), gomock.Any()).
					Return(&entity.User{
						ID:            uuid.New(),
						Name:          "test",
						Email:         "test@example.com",
						DisplayName:   "test",
						DiscordUserID: "test",
						AvatarURL:     "https://example.com/avatar.png",
					}, nil).
					Times(1)
				tp.EXPECT().
					GenerateToken(gomock.Any()).
					Return("test", nil).
					Times(1)
				tr.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&entity.Token{
						RefreshToken: "test",
					}, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "正常系: 新規ユーザーのディスコード認証成功",
			setupMock: func(m *mock.MockDiscordRepository, u *mock.MockUserRepository, tp *mock.MockTokenProvider, tr *mock.MockTokenRepository) {
				m.EXPECT().
					GetDiscordToken(gomock.Any(), gomock.Any()).
					Return(entity.DiscordToken{
						AccessToken: "test",
					}, nil).
					Times(1)
				m.EXPECT().
					FetchDiscordUser(gomock.Any(), gomock.Any()).
					Return(entity.DiscordUser{
						ID: "test",
					}, nil).
					Times(1)
				m.EXPECT().
					GetDiscordGuilds(gomock.Any(), gomock.Any()).
					Return([]string{"test"}, nil).
					Times(1)
				m.EXPECT().
					GetAllowedDiscordGuilds(gomock.Any()).
					Return([]string{"test"}, nil).
					Times(1)
				u.EXPECT().
					GetUserByDiscordUserID(gomock.Any(), gomock.Any()).
					Return(nil, domainerrors.ErrUserNotFound).
					Times(1)
				u.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&entity.User{
						ID: uuid.New(),
					}, nil).
					Times(1)
				tp.EXPECT().
					GenerateToken(gomock.Any()).
					Return("test", nil).
					Times(1)
				tr.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&entity.Token{
						RefreshToken: "test",
					}, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "異常系: ディスコード認証失敗",
			setupMock: func(m *mock.MockDiscordRepository, u *mock.MockUserRepository, tp *mock.MockTokenProvider, tr *mock.MockTokenRepository) {
				m.EXPECT().
					GetDiscordToken(gomock.Any(), gomock.Any()).
					Return(entity.DiscordToken{}, domainerrors.ErrFaileRequestToDiscord).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "異常系: ユーザーが許可されたDiscordギルドに所属していない",
			setupMock: func(m *mock.MockDiscordRepository, u *mock.MockUserRepository, tp *mock.MockTokenProvider, tr *mock.MockTokenRepository) {
				m.EXPECT().
					GetDiscordToken(gomock.Any(), gomock.Any()).
					Return(entity.DiscordToken{
						AccessToken: "test",
					}, nil).
					Times(1)
				m.EXPECT().
					FetchDiscordUser(gomock.Any(), gomock.Any()).
					Return(entity.DiscordUser{
						ID: "test",
					}, nil).
					Times(1)
				m.EXPECT().
					GetDiscordGuilds(gomock.Any(), gomock.Any()).
					Return([]string{"guild"}, nil).
					Times(1)
				m.EXPECT().
					GetAllowedDiscordGuilds(gomock.Any()).
					Return([]string{"allowed_guild"}, nil).
					Times(1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDiscordRepository := mock.NewMockDiscordRepository(ctrl)
			mockUserRepository := mock.NewMockUserRepository(ctrl)
			mockTokenProvider := mock.NewMockTokenProvider(ctrl)
			mockTokenRepository := mock.NewMockTokenRepository(ctrl)
			tt.setupMock(mockDiscordRepository, mockUserRepository, mockTokenProvider, mockTokenRepository)
			uc := NewAuthUsecase(mockDiscordRepository, mockUserRepository, mockTokenProvider, mockTokenRepository)
			appToken, refreshToken, err := uc.AuthenticateUser(context.Background(), "test")
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, appToken)
				assert.Empty(t, refreshToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, appToken)
				assert.NotEmpty(t, refreshToken)
			}
		})
	}
}
