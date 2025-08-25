package repository

import (
	"context"
	"newTiktoken/internal/video_favorite/domain/model"
)

// VideoFavoriteRepository defines the persistence interface for user favorite data.
type VideoFavoriteRepository interface {
	// Add records a user's like for a video.
	Add(ctx context.Context, favorite *model.VideoFavorite) error
	// Remove deletes a user's like for a video.
	Remove(ctx context.Context, favorite *model.VideoFavorite) error
	// ListVideoIDsByUserID returns a list of video IDs that a user has favorited.
	ListVideoIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error)
}
