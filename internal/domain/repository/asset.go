package repository

import (
	"context"
	"mime/multipart"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type AssetRepository interface {
	Create(ctx context.Context, asset *entity.Asset) (*entity.Asset, error)
	UploadFile(ctx context.Context, file *multipart.FileHeader) (*string, error)
}
