package usecase_test

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/usecase"
	"github.com/simesaba80/toybox-back/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAssetUseCase_UploadFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *mock.MockAssetRepository, file *multipart.FileHeader, userID uuid.UUID)
		wantErr bool
	}{
		{
			name: "正常系: ファイルアップロード成功",
			setup: func(t *testing.T, repo *mock.MockAssetRepository, file *multipart.FileHeader, userID uuid.UUID) {
				t.Helper()

				assetURL := "https://example.com/assets/" + uuid.NewString() + ".png"
				assetType := "image"
				var capturedAssetID uuid.UUID

				repo.EXPECT().
					UploadFile(gomock.Any(), file, gomock.AssignableToTypeOf(uuid.UUID{}), "png").
					DoAndReturn(func(ctx context.Context, fh *multipart.FileHeader, assetUUID uuid.UUID, extension string) (*string, *string, error) {
						assert.Equal(t, file, fh)
						assert.Equal(t, "png", extension)
						capturedAssetID = assetUUID
						return &assetURL, &assetType, nil
					}).
					Times(1)

				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, asset *entity.Asset) (*entity.Asset, error) {
						assert.Equal(t, capturedAssetID, asset.ID)
						assert.Equal(t, userID, asset.UserID)
						assert.Equal(t, "png", asset.Extension)
						assert.Equal(t, assetURL, asset.URL)
						assert.Equal(t, assetType, asset.AssetType)
						return asset, nil
					}).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "異常系: UploadFile でエラー",
			setup: func(t *testing.T, repo *mock.MockAssetRepository, file *multipart.FileHeader, userID uuid.UUID) {
				t.Helper()

				repo.EXPECT().
					UploadFile(gomock.Any(), file, gomock.AssignableToTypeOf(uuid.UUID{}), "png").
					Return(nil, nil, errors.New("upload failed")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "異常系: Create でエラー",
			setup: func(t *testing.T, repo *mock.MockAssetRepository, file *multipart.FileHeader, userID uuid.UUID) {
				t.Helper()

				assetURL := "https://example.com/assets/" + uuid.NewString() + ".png"
				assetType := "image"

				repo.EXPECT().
					UploadFile(gomock.Any(), file, gomock.AssignableToTypeOf(uuid.UUID{}), "png").
					Return(&assetURL, &assetType, nil).
					Times(1)

				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("create failed")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			file := &multipart.FileHeader{Filename: "test.png"}
			userID := uuid.New()

			mockRepo := mock.NewMockAssetRepository(ctrl)
			tt.setup(t, mockRepo, file, userID)

			uc := usecase.NewAssetUseCase(mockRepo)

			got, err := uc.UploadFile(context.Background(), file, userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, userID, got.UserID)
			assert.Equal(t, "png", got.Extension)
			assert.NotEmpty(t, got.AssetType)
			assert.NotEmpty(t, got.URL)
		})
	}
}
