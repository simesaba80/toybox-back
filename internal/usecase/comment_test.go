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

func TestCommentUsecase_GetCommentsByWorkID(t *testing.T) {
	tests := []struct {
		name      string
		workID    uuid.UUID
		setupMock func(*mock.MockCommentRepository, uuid.UUID)
		wantCount int
		wantErr   bool
	}{
		{
			name:   "正常系: コメント取得成功",
			workID: uuid.New(),
			setupMock: func(m *mock.MockCommentRepository, workID uuid.UUID) {
				expectedComments := []*entity.Comment{
					{
						ID:        uuid.New(),
						Content:   "Great work!",
						WorkID:    workID,
						UserID:    uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Content:   "Nice!",
						WorkID:    workID,
						UserID:    uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				m.EXPECT().
					FindByWorkID(gomock.Any(), gomock.Eq(workID)).
					Return(expectedComments, nil).
					Times(1)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:   "正常系: コメントが0件",
			workID: uuid.New(),
			setupMock: func(m *mock.MockCommentRepository, workID uuid.UUID) {
				m.EXPECT().
					FindByWorkID(gomock.Any(), gomock.Eq(workID)).
					Return([]*entity.Comment{}, nil).
					Times(1)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "異常系: リポジトリエラー",
			workID: uuid.New(),
			setupMock: func(m *mock.MockCommentRepository, workID uuid.UUID) {
				m.EXPECT().
					FindByWorkID(gomock.Any(), gomock.Eq(workID)).
					Return(nil, errors.New("database connection failed")).
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

			mockRepo := mock.NewMockCommentRepository(ctrl)
			tt.setupMock(mockRepo, tt.workID)
			mockWorkRepo := mock.NewMockWorkRepository(ctrl)
			uc := usecase.NewCommentUsecase(mockRepo, mockWorkRepo, 30*time.Second)
			got, err := uc.GetCommentsByWorkID(context.Background(), tt.workID)

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
