package usecase

import (
	"context"
	"time"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type IUserUseCase interface {
	CreateUser(ctx context.Context, name, email, passwordHash, displayName, avatar_url string) (*entity.User, error)
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

func (u *userUseCase) CreateUser(ctx context.Context, name, email, passwordHash, displayName, avatar_url string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user := &entity.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		DisplayName:  displayName,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user, err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) GetAllUser(ctx context.Context) ([]*entity.User, error) {
	return u.repo.GetAll(ctx)
}
