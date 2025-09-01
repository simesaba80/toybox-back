package user

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/simesaba80/toybox-back/internal/domain/entity"
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
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
