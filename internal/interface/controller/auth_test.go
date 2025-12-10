package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/internal/interface/controller/mock"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/pkg/echovalidator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthController_GetDiscordAuthURL(t *testing.T) {
	ctx := context.Background()

	successAuthURL := "https://discord.com/oauth2/authorize"
	successResponseBytes, _ := json.Marshal(schema.ToGetDiscordAuthURLResponse(successAuthURL))
	clientIDNotSetResponseBytes, _ := json.Marshal(map[string]string{"message": "DiscordクライアントIDが設定されていません"})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "Internal server error"})

	tests := []struct {
		name       string
		setupMock  func(mockAuthUsecase *mock.MockIAuthUsecase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "正常系: 認証URL取得成功",
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					GetDiscordAuthURL(gomock.Any()).
					DoAndReturn(func(context.Context) (string, error) {
						return successAuthURL, nil
					})
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name: "異常系: DiscordクライアントID未設定",
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					GetDiscordAuthURL(gomock.Any()).
					Return("", domainerrors.ErrClientIDNotSet)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   clientIDNotSetResponseBytes,
		},
		{
			name: "異常系: 予期しないエラー",
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					GetDiscordAuthURL(gomock.Any()).
					Return("", errors.New("unexpected error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   internalErrorResponseBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mock.NewMockIAuthUsecase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockAuthUsecase)
			}

			authController := controller.NewAuthController(mockAuthUsecase)
			e.GET("/auth/discord/", authController.GetDiscordAuthURL)

			req := httptest.NewRequest(http.MethodGet, "/auth/discord/", nil)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}

func TestAuthController_AuthenticateUser(t *testing.T) {
	successResponseBytes, _ := json.Marshal(schema.ToGetDiscordTokenResponse("app-token"))
	codeRequiredResponseBytes, _ := json.Marshal(map[string]string{"message": "code is required"})
	userNotAllowedResponseBytes, _ := json.Marshal(map[string]string{"message": "ユーザーは許可されたDiscordギルドに所属していません"})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "Internal server error"})
	failedRequestResponseBytes, _ := json.Marshal(map[string]string{"message": "Discordへのリクエストに失敗しました"})

	tests := []struct {
		name       string
		query      string
		setupMock  func(mockAuthUsecase *mock.MockIAuthUsecase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name:  "正常系: ユーザー認証成功",
			query: "?code=test-code",
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					AuthenticateUser(gomock.Any(), "test-code").
					Return("app-token", "refresh-token", nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name:       "異常系: codeが空",
			query:      "",
			setupMock:  nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   codeRequiredResponseBytes,
		},
		{
			name:  "異常系: 許可されていないギルドのユーザー",
			query: "?code=another-code",
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					AuthenticateUser(gomock.Any(), "another-code").
					Return("", "", domainerrors.ErrUserNotAllowedGuild)
			},
			wantStatus: http.StatusForbidden,
			wantBody:   userNotAllowedResponseBytes,
		},
		{
			name:  "異常系: Discordリクエスト失敗",
			query: "?code=discord-error",
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					AuthenticateUser(gomock.Any(), "discord-error").
					Return("", "", domainerrors.ErrFaileRequestToDiscord)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   failedRequestResponseBytes,
		},
		{
			name:  "異常系: 予期しないエラー",
			query: "?code=unexpected",
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					AuthenticateUser(gomock.Any(), "unexpected").
					Return("", "", errors.New("unexpected error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   internalErrorResponseBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mock.NewMockIAuthUsecase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockAuthUsecase)
			}

			authController := controller.NewAuthController(mockAuthUsecase)
			e.GET("/auth/discord/callback", authController.AuthenticateUser)

			req := httptest.NewRequest(http.MethodGet, "/auth/discord/callback"+tt.query, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}

func TestAuthController_RegenerateToken(t *testing.T) {
	successResponseBytes, _ := json.Marshal(schema.ToRegenerateTokenResponse("new-app-token"))
	bindErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "Invalid request body"})
	validationErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "Key: 'RegenerateTokenInput.RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag"})
	expiredResponseBytes, _ := json.Marshal(map[string]string{"message": "リフレッシュトークンが期限切れです"})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "Internal server error"})

	tests := []struct {
		name       string
		body       []byte
		setupMock  func(mockAuthUsecase *mock.MockIAuthUsecase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "正常系: トークン再発行成功",
			body: []byte(`{"refresh_token":"valid-refresh-token"}`),
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					RegenerateToken(gomock.Any(), "valid-refresh-token").
					Return("new-app-token", nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name:       "異常系: Bindエラー",
			body:       []byte("invalid json"),
			setupMock:  nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   bindErrorResponseBytes,
		},
		{
			name:       "異常系: バリデーションエラー",
			body:       []byte(`{"refresh_token":""}`),
			setupMock:  nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   validationErrorResponseBytes,
		},
		{
			name: "異常系: リフレッシュトークン期限切れ",
			body: []byte(`{"refresh_token":"expired-token"}`),
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					RegenerateToken(gomock.Any(), "expired-token").
					Return("", domainerrors.ErrRefreshTokenExpired)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   expiredResponseBytes,
		},
		{
			name: "異常系: 予期しないエラー",
			body: []byte(`{"refresh_token":"unexpected-token"}`),
			setupMock: func(mockAuthUsecase *mock.MockIAuthUsecase) {
				mockAuthUsecase.EXPECT().
					RegenerateToken(gomock.Any(), "unexpected-token").
					Return("", errors.New("unexpected error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   internalErrorResponseBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = echovalidator.NewValidator()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mock.NewMockIAuthUsecase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockAuthUsecase)
			}

			authController := controller.NewAuthController(mockAuthUsecase)
			e.POST("/auth/refresh", authController.RegenerateToken)

			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}
