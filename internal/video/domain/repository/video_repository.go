package repository

import (
	"context"
	"newTiktoken/internal/video/domain/model"
)

// VideoRepository defines the persistence interface for videos.
// This abstracts the data layer from the application layer.
type VideoRepository interface {
	// Feed returns a list of videos based on the latest time.
	Feed(ctx context.Context, latestTime int64) ([]*model.Video, int64, error)
	// Create saves a new video record.
	Create(ctx context.Context, video *model.Video) error
	// ListByAuthorID returns a list of videos for a given author.
	ListByAuthorID(ctx context.Context, authorID uint64) ([]*model.Video, error)
}
