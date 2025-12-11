package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type IAssetUseCase interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (*entity.Asset, error)
}

type assetUseCase struct {
	assetRepo repository.AssetRepository
}

func NewAssetUseCase(assetRepo repository.AssetRepository) IAssetUseCase {
	return &assetUseCase{
		assetRepo: assetRepo,
	}
}

func (uc *assetUseCase) UploadFile(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (*entity.Asset, error) {
	extension := strings.Split(file.Filename, ".")[1]
	asset := entity.NewAsset("", userID, extension, "")

	assetURL, assetType, err := uc.assetRepo.UploadFile(ctx, file, asset.ID, extension)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	asset.URL = *assetURL
	asset.AssetType = *assetType

	createdAsset, err := uc.assetRepo.Create(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}
	return createdAsset, nil
}
