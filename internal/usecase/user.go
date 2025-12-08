package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type IUserUseCase interface {
	GetAllUser(ctx context.Context) ([]*entity.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) IUserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) GetAllUser(ctx context.Context) ([]*entity.User, error) {
	return u.repo.GetAll(ctx)
}

func (u *userUseCase) GetByUserID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %s: %w", id.String(), err)
	}
	return user, nil
}
