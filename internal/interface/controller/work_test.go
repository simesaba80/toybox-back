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
	mockWork := &entity.Work{ID: uuid.New(), Title: "Test Work"}
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
		setupMock   func(mockWorkUsecase *mock.MockIWorkUseCase)
		wantStatus  int
		wantBody    []byte
	}{
		{
			name:        "正常系",
			queryParams: "?limit=20&page=1",
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					GetAll(gomock.Any(), util.IntPtr(20), util.IntPtr(1)).
					Return([]*entity.Work{mockWork}, 1, 20, 1, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name:        "異常系: Usecaseエラー",
			queryParams: "",
			setupMock: func(mockWorkUsecase *mock.MockIWorkUseCase) {
				mockWorkUsecase.EXPECT().
					GetAll(gomock.Any(), nil, nil).
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
			tt.setupMock(mockWorkUsecase)

			workController := controller.NewWorkController(mockWorkUsecase)
			e.GET("/works", workController.GetAllWorks)

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
	mockWork := &entity.Work{ID: workID, Title: "Test Work"}
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

func TestWorkController_CreateWork(t *testing.T) {
	userID := uuid.New()
	input := &schema.CreateWorkInput{
		Title:       "New Work",
		Description: "New Description",
		UserID:      userID.String(),
		Visibility:  "public",
	}
	inputJSON, _ := json.Marshal(input)

	createdWork := &entity.Work{ID: uuid.New(), Title: input.Title}
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
					CreateWork(gomock.Any(), input.Title, input.Description, input.Visibility, userID).
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
					CreateWork(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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

			workController := controller.NewWorkController(mockWorkUsecase)
			e.POST("/works", workController.CreateWork)

			req := httptest.NewRequest(http.MethodPost, "/works", bytes.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}
