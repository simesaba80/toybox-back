package controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/internal/interface/controller/mock"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserController_GetAllUsers(t *testing.T) {
	mockUser := &entity.User{ID: uuid.New(), Name: "testuser"}
	successResponseBytes, _ := json.Marshal(schema.UserListResponse{
		Users: []schema.GetUserOutput{schema.ToUserResponse(mockUser)},
	})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "サーバーエラーが発生しました"})

	tests := []struct {
		name       string
		setupMock  func(mockUserUsecase *mock.MockIUserUseCase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "正常系",
			setupMock: func(mockUserUsecase *mock.MockIUserUseCase) {
				mockUserUsecase.EXPECT().
					GetAllUser(gomock.Any()).
					Return([]*entity.User{mockUser}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name: "異常系: Usecaseエラー",
			setupMock: func(mockUserUsecase *mock.MockIUserUseCase) {
				mockUserUsecase.EXPECT().
					GetAllUser(gomock.Any()).
					Return(nil, errors.New("some error"))
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

			mockUsecase := mock.NewMockIUserUseCase(ctrl)
			tt.setupMock(mockUsecase)

			userController := controller.NewUserController(mockUsecase)
			e.GET("/users", userController.GetAllUsers)

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}

func TestUserController_GetUserByID(t *testing.T) {
	mockUser := &entity.User{ID: uuid.New(), Name: "testuser"}
	successResponseBytes, _ := json.Marshal(schema.ToUserResponse(mockUser))
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "サーバーエラーが発生しました"})

	tests := []struct {
		name       string
		setupMock  func(mockUserUsecase *mock.MockIUserUseCase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "正常系",
			setupMock: func(mockUserUsecase *mock.MockIUserUseCase) {
				mockUserUsecase.EXPECT().
					GetByUserID(gomock.Any(), gomock.Eq(mockUser.ID)).
					Return(mockUser, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name: "異常系: Usecaseエラー",
			setupMock: func(mockUserUsecase *mock.MockIUserUseCase) {
				mockUserUsecase.EXPECT().
					GetByUserID(gomock.Any(), gomock.Eq(mockUser.ID)).
					Return(nil, errors.New("some error"))
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

			mockUsecase := mock.NewMockIUserUseCase(ctrl)
			tt.setupMock(mockUsecase)

			userController := controller.NewUserController(mockUsecase)
			e.GET("/users/:id", userController.GetUserByID)

			req := httptest.NewRequest(http.MethodGet, "/users/"+mockUser.ID.String(), nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}
