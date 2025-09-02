package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type WorkUseCase struct {
	repo repository.WorkRepository
}

func NewWorkUseCase(repo repository.WorkRepository) *WorkUseCase {
	return &WorkUseCase{
		repo: repo,
	}
}

func (uc *WorkUseCase) GetAll(ctx context.Context) ([]*entity.Work, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *WorkUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	return uc.repo.GetByID(ctx, id)
}
