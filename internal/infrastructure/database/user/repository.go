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
	dtoUser := dto.ToUserDTO(user)

	_, err := r.db.NewInsert().Model(dtoUser).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return dtoUser.ToUserEntity(), nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	dtoUsers := make([]*dto.User, 0)
	err := r.db.NewSelect().Model(&dtoUsers).Scan(ctx)
	if err != nil {
		return nil, err
	}

	entityUsers := make([]*entity.User, len(dtoUsers))
	for i, dtoUser := range dtoUsers {
		entityUsers[i] = dtoUser.ToUserEntity()
	}

	return entityUsers, nil
}

func (r *UserRepository) GetUserByDiscordUserID(ctx context.Context, discordUserID string) (*entity.User, error) {
	dtoUser := new(dto.User)
	err := r.db.NewSelect().Model(dtoUser).Where("discord_user_id = ?", discordUserID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return dtoUser.ToUserEntity(), nil
}
