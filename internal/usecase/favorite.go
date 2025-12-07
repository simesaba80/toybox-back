package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type IFavoriteUsecase interface {
	CreateFavorite(ctx context.Context, workID uuid.UUID, userID uuid.UUID) (*entity.Favorite, error)
	DeleteFavorite(ctx context.Context, workID uuid.UUID, userID uuid.UUID) error
	CountFavoritesByWorkID(ctx context.Context, workID uuid.UUID) (int, error)
	IsFavorite(ctx context.Context, workID uuid.UUID, userID uuid.UUID) (bool, error)
}

type favoriteUsecase struct {
	favoriteRepo repository.FavoriteRepository
}

func NewFavoriteUsecase(favoriteRepo repository.FavoriteRepository) IFavoriteUsecase {
	return &favoriteUsecase{
		favoriteRepo: favoriteRepo,
	}
}

func (uc *favoriteUsecase) CreateFavorite(ctx context.Context, workID uuid.UUID, userID uuid.UUID) (*entity.Favorite, error) {
	favorite := entity.NewFavorite(workID, userID)
	exists, err := uc.favoriteRepo.Exists(ctx, favorite)
	if err != nil {
		return nil, fmt.Errorf("failed to check if favorite exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("favorite already exists")
	}

	createdFavorite, err := uc.favoriteRepo.Create(ctx, favorite)
	if err != nil {
		return nil, fmt.Errorf("failed to create favorite: %w", err)
	}
	return createdFavorite, nil
}

func (uc *favoriteUsecase) DeleteFavorite(ctx context.Context, workID uuid.UUID, userID uuid.UUID) error {
	favorite := entity.NewFavorite(workID, userID)
	exists, err := uc.favoriteRepo.Exists(ctx, favorite)
	if err != nil {
		return fmt.Errorf("failed to check if favorite exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("favorite not found")
	}
	return uc.favoriteRepo.Delete(ctx, favorite)
}

func (uc *favoriteUsecase) CountFavoritesByWorkID(ctx context.Context, workID uuid.UUID) (int, error) {
	return uc.favoriteRepo.CountByWorkID(ctx, workID)
}

func (uc *favoriteUsecase) IsFavorite(ctx context.Context, workID uuid.UUID, userID uuid.UUID) (bool, error) {
	favorite := entity.NewFavorite(workID, userID)
	return uc.favoriteRepo.Exists(ctx, favorite)
}
