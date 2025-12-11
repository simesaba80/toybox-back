package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type IWorkUseCase interface {
	GetAll(ctx context.Context, limit, page *int, userID uuid.UUID) ([]*entity.Work, int, int, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, authenticatedUserID uuid.UUID) ([]*entity.Work, error)
	CreateWork(ctx context.Context, title, description, visibility string, thumbnailAssetID uuid.UUID, assetIDs []uuid.UUID, urls []string, userID uuid.UUID) (*entity.Work, error)
}

type workUseCase struct {
	repo repository.WorkRepository
}

func NewWorkUseCase(repo repository.WorkRepository) IWorkUseCase {
	return &workUseCase{
		repo: repo,
	}
}

func (uc *workUseCase) GetAll(ctx context.Context, limit, page *int, userID uuid.UUID) ([]*entity.Work, int, int, int, error) {
	actualLimit := 20
	actualPage := 1
	if limit != nil {
		actualLimit = *limit
	}
	if page != nil {
		actualPage = *page
	}
	fmt.Println("actualLimit", actualLimit)
	fmt.Println("actualPage", actualPage)
	fmt.Println("userID", userID)

	offset := (actualPage - 1) * actualLimit
	if userID == uuid.Nil {
		works, total, err := uc.repo.GetAllPublic(ctx, actualLimit, offset)
		if err != nil {
			return nil, 0, 0, 0, fmt.Errorf("failed to get all works by user ID %s: %w", userID.String(), err)
		}
		return works, total, actualLimit, actualPage, nil
	}

	works, total, err := uc.repo.GetAll(ctx, actualLimit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to get all works: %w", err)
	}
	return works, total, actualLimit, actualPage, nil
}

func (uc *workUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	work, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work by ID %s: %w", id.String(), err)
	}
	return work, nil
}

func (uc *workUseCase) GetByUserID(ctx context.Context, userID uuid.UUID, authenticatedUserID uuid.UUID) ([]*entity.Work, error) {
	var public bool
	if authenticatedUserID == uuid.Nil {
		public = true
	} else {
		public = false
	}

	works, err := uc.repo.GetByUserID(ctx, userID, public)
	if err != nil {
		return nil, fmt.Errorf("failed to get works by user ID %s: %w", userID.String(), err)
	}
	return works, nil
}

func (uc *workUseCase) CreateWork(ctx context.Context, title, description, visibility string, thumbnailAssetID uuid.UUID, assetIDs []uuid.UUID, urls []string, userID uuid.UUID) (*entity.Work, error) {
	if title == "" {
		return nil, domainerrors.ErrInvalidTitle
	}
	if description == "" {
		return nil, domainerrors.ErrInvalidDescription
	}
	if visibility == "" {
		return nil, domainerrors.ErrInvalidVisibility
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

	work := entity.NewWork(title, description, userID, visibility, thumbnailAssetID, assets, urlPointers)

	createdWork, err := uc.repo.Create(ctx, work)
	if err != nil {
		return nil, fmt.Errorf("failed to create work: %w", err)
	}
	return createdWork, nil
}
