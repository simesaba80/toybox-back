package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type WorkUseCase struct {
	repo    repository.WorkRepository
	timeout time.Duration
}

func NewWorkUseCase(repo repository.WorkRepository, timeout time.Duration) *WorkUseCase {
	return &WorkUseCase{
		repo:    repo,
		timeout: time.Second * 30,
	}
}

func (uc *WorkUseCase) GetAll(ctx context.Context) ([]*entity.Work, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	works, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all works: %w", err)
	}
	return works, nil
}

func (uc *WorkUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	work, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work by ID %s: %w", id.String(), err)
	}
	return work, nil
}

func (uc *WorkUseCase) CreateWork(ctx context.Context, title, description, descriptionHTML, visibility string, userID uuid.UUID) (*entity.Work, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	work := &entity.Work{
		ID:              uuid.New(),
		Title:           title,
		Description:     description,
		DescriptionHTML: descriptionHTML,
		UserID:          userID,
		Visibility:      visibility,
	}

	if err := work.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	createdWork, err := uc.repo.Create(ctx, work)
	if err != nil {
		return nil, fmt.Errorf("failed to create work: %w", err)
	}
	return createdWork, nil
}
