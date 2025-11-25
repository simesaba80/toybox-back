package usecase

import (
	"context"
	"time"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type IUserUseCase interface {
	GetAllUser(ctx context.Context) ([]*entity.User, error)
}

type userUseCase struct {
	repo    repository.UserRepository
	timeout time.Duration
}

func NewUserUseCase(repo repository.UserRepository, timeout time.Duration) IUserUseCase {
	return &userUseCase{
		repo:    repo,
		timeout: time.Second * 30,
	}
}

func (u *userUseCase) GetAllUser(ctx context.Context) ([]*entity.User, error) {
	return u.repo.GetAll(ctx)
}
