package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type IWorkUseCase interface {
	GetAll(ctx context.Context, limit, page *int) ([]*entity.Work, int, int, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error)
	CreateWork(ctx context.Context, title, description, visibility string, thumbnailAssetID uuid.UUID, assetIDs []uuid.UUID, urls []string, userID uuid.UUID, tagIDs []uuid.UUID) (*entity.Work, error)
}

type workUseCase struct {
	workRepo repository.WorkRepository
	tagRepo  repository.TagRepository
	timeout  time.Duration
}

func NewWorkUseCase(workRepo repository.WorkRepository, tagRepo repository.TagRepository, timeout time.Duration) IWorkUseCase {
	return &workUseCase{
		workRepo: workRepo,
		tagRepo:  tagRepo,
		timeout:  time.Second * 30,
	}
}

func (uc *workUseCase) GetAll(ctx context.Context, limit, page *int) ([]*entity.Work, int, int, int, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	actualLimit := 20
	actualPage := 1
	if limit != nil {
		actualLimit = *limit
	}
	if page != nil {
		actualPage = *page
	}

	offset := (actualPage - 1) * actualLimit

	works, total, err := uc.workRepo.GetAll(ctx, actualLimit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to get all works: %w", err)
	}
	return works, total, actualLimit, actualPage, nil
}

func (uc *workUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	work, err := uc.workRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work by ID %s: %w", id.String(), err)
	}
	return work, nil
}

func (uc *workUseCase) CreateWork(ctx context.Context, title, description, visibility string, thumbnailAssetID uuid.UUID, assetIDs []uuid.UUID, urls []string, userID uuid.UUID, tagIDs []uuid.UUID) (*entity.Work, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()
	if title == "" {
		return nil, domainerrors.ErrInvalidTitle
	}
	if description == "" {
		return nil, domainerrors.ErrInvalidDescription
	}
	if visibility == "" {
		return nil, domainerrors.ErrInvalidVisibility
	}
	if len(tagIDs) == 0 {
		return nil, domainerrors.ErrInvalidTagIDs
	}

	var tags []*entity.Tag
	var err error

	exists, err := uc.tagRepo.ExistAll(ctx, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to check tag existence: %w", err)
	}
	if !exists {
		return nil, domainerrors.ErrTagNotFound
	}

	tags, err = uc.tagRepo.FindAllByIDs(ctx, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find tags by ids: %w", err)
	}

	assets := make([]*entity.Asset, len(assetIDs))
	for i, assetID := range assetIDs {
		assets[i] = &entity.Asset{
			ID: assetID,
		}
	}

	urlPointers := make([]*string, len(urls))
	for i, url := range urls {
		urlPointers[i] = &url
	}

	work := entity.NewWork(title, description, userID, visibility, thumbnailAssetID, assets, urlPointers, tagIDs, tags)

	createdWork, err := uc.workRepo.Create(ctx, work)
	if err != nil {
		return nil, fmt.Errorf("failed to create work: %w", err)
	}
	return createdWork, nil
}
