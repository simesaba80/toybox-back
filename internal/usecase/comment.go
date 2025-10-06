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

func (uc *CommentUsecase) CreateComment(ctx context.Context, content string, workID, userID uuid.UUID, replyAt string) (*entity.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	// replyAtがある場合は返信先にコメントが存在するか確認
	if replyAt != "" {
		replyID, err := uuid.Parse(replyAt)
		if err != nil {
			return nil, fmt.Errorf("invalid reply_at format: %w", err)
		}
		_, err = uc.repo.FindByID(ctx, replyID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate reply target comment %s: %w", replyAt, err)
		}
	}

	comment := &entity.Comment{
		ID:      uuid.New(),
		Content: content,
		WorkID:  workID,
		UserID:  userID,
		ReplyAt: replyAt,
	}

	createdComment, err := uc.repo.Create(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return createdComment, nil
}
