package asset

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

type AssetRepository struct {
	db *bun.DB
	s3 *s3.Client
}

func NewAssetRepository(db *bun.DB, s3 *s3.Client) *AssetRepository {
	return &AssetRepository{
		db: db,
		s3: s3,
	}
}

func (r *AssetRepository) Create(ctx context.Context, asset *entity.Asset) (*entity.Asset, error) {
	dtoAsset := dto.ToAssetDTO(asset)
	_, err := r.db.NewInsert().Model(dtoAsset).Exec(ctx)
	if err != nil {
		return nil, domainerrors.ErrFailedToCreateAsset
	}
	return dtoAsset.ToAssetEntity(), nil
}

func (r *AssetRepository) UploadFile(ctx context.Context, file *multipart.FileHeader) (*string, error) {
	assetURL := "https://example.com/asset.png"
	return &assetURL, nil
}
