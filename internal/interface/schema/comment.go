package schema

import (
	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/domain/entity"
)

type UserInCommentResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type CommentResponse struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	ReplyAt   string                 `json:"reply_at"`
	User      *UserInCommentResponse `json:"user"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

func ToCommentResponse(comment *entity.Comment) *CommentResponse {
	if comment == nil {
		return nil
	}

	var user *UserInCommentResponse
	if comment.User != nil && comment.User.ID != uuid.Nil {
		user = &UserInCommentResponse{
			ID:          comment.User.ID.String(),
			DisplayName: comment.User.DisplayName,
			AvatarURL:   comment.User.AvatarURL,
		}
	}

	return &CommentResponse{
		ID:        comment.ID.String(),
		Content:   comment.Content,
		ReplyAt:   comment.ReplyAt,
		User:      user,
		CreatedAt: comment.CreatedAt.String(),
		UpdatedAt: comment.UpdatedAt.String(),
	}
}

func ToCommentListResponse(comments []*entity.Comment) []*CommentResponse {
	res := make([]*CommentResponse, len(comments))
	for i, comment := range comments {
		res[i] = ToCommentResponse(comment)
	}
	return res
}
