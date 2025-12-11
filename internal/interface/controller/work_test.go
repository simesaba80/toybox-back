package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/internal/interface/controller/mock"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/util"
	"github.com/simesaba80/toybox-back/pkg/echovalidator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorkController_GetAllWorks(t *testing.T) {
	userID := uuid.New()
	mockWork := &entity.Work{
		ID:        uuid.New(),
		Title:     "Test Work",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	successResponseBytes, _ := json.Marshal(schema.WorkListResponse{
		Works:      []schema.GetWorkOutput{schema.ToWorkResponse(mockWork)},
		TotalCount: 1,
		Page:       1,
		Limit:      20,
	})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "サーバーエラーが発生しました"})

	tests := []struct {
		name        string
		queryParams string
		withAuth    bool
		userID      uuid.UUID
		setupMock   func(mockWorkUsecase *mock.MockIWorkUseCase, userID uuid.UUID)
		wantStatus  int
		wantBody    []byte
	}{
		{
			name:        "正常系: 認証あり",
			queryParams: "?limit=20&page=1",
			withAuth:    true,
			userID:      userID,
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase, userID uuid.UUID) {
				mockWorkUsecase.EXPECT().
					GetAll(gomock.Any(), util.IntPtr(20), util.IntPtr(1), userID).
					Return([]*entity.Work{mockWork}, 1, 20, 1, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name:        "正常系: 認証なし（公開作品のみ）",
			queryParams: "?limit=20&page=1",
			withAuth:    false,
			userID:      uuid.Nil,
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase, userID uuid.UUID) {
				mockWorkUsecase.EXPECT().
					GetAll(gomock.Any(), util.IntPtr(20), util.IntPtr(1), uuid.Nil).
					Return([]*entity.Work{mockWork}, 1, 20, 1, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name:        "異常系: Usecaseエラー（認証なし）",
			queryParams: "",
			withAuth:    false,
			userID:      uuid.Nil,
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase, userID uuid.UUID) {
				mockWorkUsecase.EXPECT().
					GetAll(gomock.Any(), nil, nil, uuid.Nil).
					Return(nil, 0, 0, 0, errors.New("some error"))
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

			mockWorkUsecase := mock.NewMockIWorkUseCase(ctrl)
			tt.setupMock(mockWorkUsecase, tt.userID)

			workController := controller.NewWorkController(mockWorkUsecase)
			e.GET("/works", func(c echo.Context) error {
				if tt.withAuth {
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, &schema.JWTCustomClaims{
						UserID: tt.userID.String(),
					})
					c.Set("user", token)
				}
				return workController.GetAllWorks(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/works"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}

func TestWorkController_GetWorkByID(t *testing.T) {
	workID := uuid.New()
	mockWork := &entity.Work{
		ID:        workID,
		Title:     "Test Work",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	successResponseBytes, _ := json.Marshal(schema.ToWorkResponse(mockWork))
	invalidIDResponseBytes, _ := json.Marshal(map[string]string{"message": "無効なリクエストです"})
	notFoundResponseBytes, _ := json.Marshal(map[string]string{"message": "作品が見つかりませんでした"})

	tests := []struct {
		name       string
		workID     string
		setupMock  func(mockWorkUsecase *mock.MockIWorkUseCase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name:   "正常系",
			workID: workID.String(),
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					GetByID(gomock.Any(), workID).
					Return(mockWork, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name:       "異常系: work_idが不正",
			workID:     "invalid-uuid",
			setupMock:  func(mockWorkUsecase *mock.MockIWorkUseCase) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   invalidIDResponseBytes,
		},
		{
			name:   "異常系: Not Found",
			workID: workID.String(),
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					GetByID(gomock.Any(), workID).
					Return(nil, domainerrors.ErrWorkNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   notFoundResponseBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWorkUsecase := mock.NewMockIWorkUseCase(ctrl)
			tt.setupMock(mockWorkUsecase)

			workController := controller.NewWorkController(mockWorkUsecase)
			e.GET("/works/:work_id", workController.GetWorkByID)

			req := httptest.NewRequest(http.MethodGet, "/works/"+tt.workID, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}

func TestWorkController_GetWorksByUserID(t *testing.T) {
	targetUserID := uuid.New()
	authenticatedUserID := uuid.New()
	mockWork1 := &entity.Work{
		ID:        uuid.New(),
		Title:     "Public Work",
		UserID:    targetUserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockWork2 := &entity.Work{
		ID:        uuid.New(),
		Title:     "Private Work",
		UserID:    targetUserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	successResponseWithBothWorks, _ := json.Marshal(schema.ToWorkListResponse([]*entity.Work{mockWork1, mockWork2}))
	successResponseWithPublicOnly, _ := json.Marshal(schema.ToWorkListResponse([]*entity.Work{mockWork1}))
	badRequestResponseBytes, _ := json.Marshal(map[string]string{"message": "無効なリクエストボディです"})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "サーバーエラーが発生しました"})

	tests := []struct {
		name       string
		userID     string
		withAuth   bool
		authUserID uuid.UUID
		setupMock  func(mockWorkUsecase *mock.MockIWorkUseCase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "正常系: 認証あり",
			userID:     targetUserID.String(),
			withAuth:   true,
			authUserID: authenticatedUserID,
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					GetByUserID(gomock.Any(), targetUserID, authenticatedUserID).
					Return([]*entity.Work{mockWork1, mockWork2}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseWithBothWorks,
		},
		{
			name:       "正常系: 認証なし（公開作品のみ）",
			userID:     targetUserID.String(),
			withAuth:   false,
			authUserID: uuid.Nil,
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					GetByUserID(gomock.Any(), targetUserID, uuid.Nil).
					Return([]*entity.Work{mockWork1}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseWithPublicOnly,
		},
		{
			name:       "異常系: user_idが不正",
			userID:     "invalid-uuid",
			withAuth:   false,
			authUserID: uuid.Nil,
			setupMock:  func(mockWorkUsecase *mock.MockIWorkUseCase) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   badRequestResponseBytes,
		},
		{
			name:       "異常系: Usecaseエラー",
			userID:     targetUserID.String(),
			withAuth:   false,
			authUserID: uuid.Nil,
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					GetByUserID(gomock.Any(), targetUserID, uuid.Nil).
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

			mockWorkUsecase := mock.NewMockIWorkUseCase(ctrl)
			tt.setupMock(mockWorkUsecase)

			workController := controller.NewWorkController(mockWorkUsecase)
			e.GET("/works/:user_id", func(c echo.Context) error {
				if tt.withAuth {
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, &schema.JWTCustomClaims{
						UserID: tt.authUserID.String(),
					})
					c.Set("user", token)
				}
				return workController.GetWorksByUserID(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/works/"+tt.userID, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}

func TestWorkController_CreateWork(t *testing.T) {
	userID := uuid.New()
	input := &schema.CreateWorkInput{
		Title:            "New Work",
		Description:      "New Description",
		ThumbnailAssetID: uuid.New(),
		AssetIDs:         []uuid.UUID{uuid.New()},
		Visibility:       "public",
		URLs:             []string{"https://example.com"},
	}
	inputJSON, _ := json.Marshal(input)

	createdWork := &entity.Work{
		ID:        uuid.New(),
		Title:     input.Title,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	successResponseBytes, _ := json.Marshal(schema.ToCreateWorkOutput(createdWork))
	badRequestResponseBytes, _ := json.Marshal(map[string]string{"message": "無効なリクエストボディです"})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "サーバーエラーが発生しました"})

	tests := []struct {
		name       string
		body       []byte
		setupMock  func(mockWorkUsecase *mock.MockIWorkUseCase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "正常系",
			body: inputJSON,
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					CreateWork(gomock.Any(), input.Title, input.Description, input.Visibility, input.ThumbnailAssetID, input.AssetIDs, input.URLs, userID).
					Return(createdWork, nil)
			},
			wantStatus: http.StatusCreated,
			wantBody:   successResponseBytes,
		},
		{
			name:       "異常系: 不正なリクエストボディ",
			body:       []byte("invalid json"),
			setupMock:  func(mockWorkUsecase *mock.MockIWorkUseCase) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   badRequestResponseBytes,
		},
		{
			name: "異常系: Usecaseエラー",
			body: inputJSON,
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					CreateWork(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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

			mockWorkUsecase := mock.NewMockIWorkUseCase(ctrl)
			tt.setupMock(mockWorkUsecase)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, &schema.JWTCustomClaims{
				UserID: userID.String(),
			})

			workController := controller.NewWorkController(mockWorkUsecase)
			e.POST("/works", func(c echo.Context) error {
				c.Set("user", token)
				return workController.CreateWork(c)
			})

			req := httptest.NewRequest(http.MethodPost, "/works", bytes.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}
