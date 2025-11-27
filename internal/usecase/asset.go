package usecase

import (
	"context"
	"fmt"
	"mime/multipart"

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
	assetURL, err := uc.assetRepo.UploadFile(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	// asset := &entity.Asset{
	// 	AssetType: "",
	// 	UserID:    userID,
	// 	Extension: strings.Split(file.Filename, ".")[1],
	// 	URL:       *assetURL,
	// }
	// createdAsset, err := uc.assetRepo.Create(ctx, asset)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create asset: %w", err)
	// }
	fmt.Println("UserID: ", userID)
	return assetURL, nil
}
