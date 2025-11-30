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
	CreateWork(ctx context.Context, title, description, visibility, thumbnailAssetID string, assetIDs []string, userID string) (*entity.Work, error)
}

type workUseCase struct {
	repo    repository.WorkRepository
	timeout time.Duration
}

func NewWorkUseCase(repo repository.WorkRepository, timeout time.Duration) IWorkUseCase {
	return &workUseCase{
		repo:    repo,
		timeout: time.Second * 30,
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

	works, total, err := uc.repo.GetAll(ctx, actualLimit, offset)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to get all works: %w", err)
	}
	return works, total, actualLimit, actualPage, nil
}

func (uc *workUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	work, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work by ID %s: %w", id.String(), err)
	}
	return work, nil
}

func (uc *workUseCase) CreateWork(ctx context.Context, title, description, visibility, thumbnailAssetID string, assetIDs []string, userID string) (*entity.Work, error) {
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

	work := &entity.Work{
		ID:          "",
		Title:       title,
		Description: description,
		UserID:      userID,
		Visibility:  visibility,
	}

	createdWork, err := uc.repo.Create(ctx, work)
	if err != nil {
		return nil, fmt.Errorf("failed to create work: %w", err)
	}
	return createdWork, nil
}
