package asset

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
	"github.com/uptrace/bun"
)

type mimeType string

const (
	png             mimeType = "image/png"
	jpeg            mimeType = "image/jpeg"
	jpg             mimeType = "image/jpge"
	bmp             mimeType = "image/bmp"
	gif             mimeType = "image/gif"
	webp            mimeType = "image/webp"
	mp4             mimeType = "video/mp4"
	mov             mimeType = "video/quicktime"
	avi             mimeType = "video/x-msvideo"
	flv             mimeType = "video/x-flv"
	mp3             mimeType = "audio/mpeg"
	wav             mimeType = "audio/wav"
	m4a             mimeType = "audio/aac"
	zip             mimeType = "application/zip"
	gilf            mimeType = "model/gltf+json"
	defaultMimeType mimeType = "application/octet-stream"
)

var ExtensionToDirName = map[string]string{
	"png":  "image",
	"jpeg": "image",
	"jpg":  "image",
	"bmp":  "image",
	"gif":  "image",
	"webp": "image",
	"mp4":  "video",
	"mov":  "video",
	"avi":  "video",
	"flv":  "video",
	"mp3":  "music",
	"wav":  "music",
	"m4a":  "music",
	"zip":  "zip",
	"gltf": "model",
	"fbx":  "model",
}

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
	openFile, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer openFile.Close()
	extension := strings.Split(file.Filename, ".")[1]
	mimeType := mimeType(extension)
	if mimeType == "" {
		mimeType = defaultMimeType
	}
	dirName := ExtensionToDirName[extension]
	if dirName == "" {
		dirName = "other"
	}
	assetUUID := uuid.New().String()
	_, err = r.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(config.S3_BUCKET),
		Key:         aws.String(config.S3_DIR + "/" + dirName + "/" + assetUUID + "/origin." + extension),
		Body:        openFile,
		ContentType: aws.String(string(mimeType)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	assetURL := config.S3_BASE_URL + "/" + config.S3_BUCKET + "/" + config.S3_DIR + "/" + dirName + "/" + assetUUID + "/origin." + extension
	return &assetURL, nil
}
