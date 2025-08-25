package application

import (
	"context"
	"log/slog"
	"newTiktoken/internal/video_comment/domain/model"
	"newTiktoken/internal/video_comment/domain/repository"
	"newTiktoken/pkg/pb/video_comment"
)

// CommentApplicationService provides comment-related application services.
type CommentApplicationService struct {
	commentRepo repository.VideoCommentRepository
	logger      *slog.Logger
}

// NewCommentApplicationService creates a new CommentApplicationService.
func NewCommentApplicationService(commentRepo repository.VideoCommentRepository, logger *slog.Logger) *CommentApplicationService {
	return &CommentApplicationService{commentRepo: commentRepo, logger: logger}
}

// CommentAction handles the comment action.
func (s *CommentApplicationService) CommentAction(ctx context.Context, req *video_comment.CommentActionRequest) (*model.Comment, error) {
	// TODO: Implement business logic
	return nil, nil
}

// CommentList handles the comment list query.
func (s *CommentApplicationService) CommentList(ctx context.Context, req *video_comment.CommentListRequest) ([]*model.Comment, error) {
	// TODO: Implement business logic
	return nil, nil
}
