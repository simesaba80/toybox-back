package user

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/infrastructure/database/dto"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	dtoUser := toDTO(user)
	_, err := r.db.NewInsert().Model(dtoUser).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	var dtoUsers []*dto.User
	err := r.db.NewSelect().Model(&dtoUsers).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return toEntities(dtoUsers), nil
}
