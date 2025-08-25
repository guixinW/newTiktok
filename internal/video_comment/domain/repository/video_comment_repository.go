package repository

import (
	"context"
	"newTiktoken/internal/video_comment/domain/model"
)

// VideoCommentRepository defines the interface for comment data operations.
type VideoCommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, commentID uint64) error
	FindByVideoID(ctx context.Context, videoID uint64) ([]*model.Comment, error)
}
