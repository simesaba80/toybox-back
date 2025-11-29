package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type IAssetUseCase interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, userID string) (*string, error)
}

type assetUseCase struct {
	assetRepo repository.AssetRepository
}

func NewAssetUseCase(assetRepo repository.AssetRepository) IAssetUseCase {
	return &assetUseCase{
		assetRepo: assetRepo,
	}
}

func (uc *assetUseCase) UploadFile(ctx context.Context, file *multipart.FileHeader, userID string) (*string, error) {
	extension := strings.Split(file.Filename, ".")[1]

	assetURL, assetUUID, err := uc.assetRepo.UploadFile(ctx, file, extension)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	asset := &entity.Asset{
		ID:        *assetUUID,
		WorkID:    "",
		UserID:    userID,
		AssetType: "",
		Extension: extension,
		URL:       *assetURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdAsset, err := uc.assetRepo.Create(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}
	return &createdAsset.URL, nil
}
