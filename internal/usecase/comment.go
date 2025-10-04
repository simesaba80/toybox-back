package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
	"github.com/simesaba80/toybox-back/internal/domain/repository"
)

type CommentUsecase struct {
	repo    repository.CommentRepository
	timeout time.Duration
}

func NewCommentUsecase(repo repository.CommentRepository, timeout time.Duration) *CommentUsecase {
	return &CommentUsecase{
		repo:    repo,
		timeout: time.Second * 30,
	}
}

func (uc *CommentUsecase) GetCommentsByWorkID(ctx context.Context, workID uuid.UUID) ([]*entity.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	comments, err := uc.repo.FindByWorkID(ctx, workID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by work ID %s: %w", workID.String(), err)
	}

	return comments, nil
}
