package controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/internal/interface/controller/mock"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCommentController_GetCommentsByWorkID(t *testing.T) {
	workID := uuid.New()

	mockComments := []*entity.Comment{
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			WorkID:    workID,
			Content:   "コメント",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	successResponseBytes, _ := json.Marshal(schema.ToCommentListResponse(mockComments))
	badRequestResponseBytes, _ := json.Marshal(map[string]string{"message": "Invalid work ID format"})
	internalErrorResponseBytes, _ := json.Marshal(map[string]string{"message": "Failed to retrieve comments"})

	tests := []struct {
		name       string
		workID     string
		setupMock  func(mockCommentUsecase *mock.MockICommentUsecase, mockWorkUsecase *mock.MockIWorkUseCase)
		wantStatus int
		wantBody   []byte
	}{
		{
			name:   "正常系: コメント取得成功",
			workID: workID.String(),
			setupMock: func(mockCommentUsecase *mock.MockICommentUsecase, mockWorkUsecase *mock.MockIWorkUseCase) {
				mockCommentUsecase.EXPECT().
					GetCommentsByWorkID(gomock.Any(), workID).
					Return(mockComments, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   successResponseBytes,
		},
		{
			name:   "異常系: work_idが不正",
			workID: "invalid-uuid",
			setupMock: func(mockCommentUsecase *mock.MockICommentUsecase, mockWorkUsecase *mock.MockIWorkUseCase) {
				mockCommentUsecase.EXPECT().
					GetCommentsByWorkID(gomock.Any(), workID).
					Return(nil, errors.New("some db error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   badRequestResponseBytes,
		},
		{
			name:   "異常系: Usecaseエラー",
			workID: workID.String(),
			setupMock: func(mockCommentUsecase *mock.MockICommentUsecase, mockWorkUsecase *mock.MockIWorkUseCase) {
				mockCommentUsecase.EXPECT().
					GetCommentsByWorkID(gomock.Any(), workID).
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

			mockCommentUsecase := mock.NewMockICommentUsecase(ctrl)
			mockWorkUsecase := mock.NewMockIWorkUseCase(ctrl)
			tt.setupMock(mockCommentUsecase, mockWorkUsecase)

			commentController := controller.NewCommentController(mockCommentUsecase)
			e.GET("/works/:work_id/comments", commentController.GetCommentsByWorkID)

			req := httptest.NewRequest(http.MethodGet, "/works/"+tt.workID+"/comments", nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, string(tt.wantBody), rec.Body.String())
		})
	}
}
