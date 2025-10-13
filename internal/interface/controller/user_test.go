package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/pkg/echovalidator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserController_CreateUser(t *testing.T) {
	input := &schema.CreateUserInput{
		Name:         "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		DisplayName:  "Test User",
		AvatarURL:    "http://example.com/avatar.png",
	}
	inputJSON, _ := json.Marshal(input)

	createdUser := &entity.User{ID: uuid.New(), Name: input.Name, Email: input.Email}
	successResponseBytes, _ := json.Marshal(schema.ToUserResponse(createdUser))
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "Internal Server Error"})

	tests := []struct {
		name       string
		body       []byte
		setupMock  func(m *controller.MockUserUseCase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "正常系",
			body: inputJSON,
			setupMock: func(m *controller.MockUserUseCase) {
				m.EXPECT().
					CreateUser(gomock.Any(), input.Name, input.Email, input.PasswordHash, input.DisplayName, input.AvatarURL).
					Return(createdUser, nil)
			},
			wantStatus: http.StatusCreated,
			wantBody:   successResponseBytes,
		},
		{
			name: "異常系: Usecaseエラー",
			body: inputJSON,
			setupMock: func(m *controller.MockUserUseCase) {
				m.EXPECT().
					CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some error"))
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

			mockUsecase := controller.NewMockUserUseCase(ctrl)
			tt.setupMock(mockUsecase)

			userController := controller.NewUserController(mockUsecase)
			e.POST("/users", userController.CreateUser)

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}

func TestUserController_GetAllUsers(t *testing.T) {
	mockUser := &entity.User{ID: uuid.New(), Name: "testuser"}
	successResponseBytes, _ := json.Marshal(schema.UserListResponse{
		Users: []schema.GetUserOutput{schema.ToUserResponse(mockUser)},
	})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "Internal Server Error"})

	tests := []struct {
		name       string
		setupMock  func(m *controller.MockUserUseCase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "正常系",
			setupMock: func(m *controller.MockUserUseCase) {
				m.EXPECT().
					GetAllUser(gomock.Any()).
					Return([]*entity.User{mockUser}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name: "異常系: Usecaseエラー",
			setupMock: func(m *controller.MockUserUseCase) {
				m.EXPECT().
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

			mockUsecase := controller.NewMockUserUseCase(ctrl)
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
