package persistence

import (
	"context"
	"newTiktoken/internal/video_favorite/domain/model"
	"sync"
)

// InMemoryVideoFavoriteRepository is an in-memory implementation of VideoFavoriteRepository.
type InMemoryVideoFavoriteRepository struct {
	mu        sync.RWMutex
	favorites map[uint64]map[uint64]struct{} // userID -> {videoID -> {}}
}

// NewInMemoryVideoFavoriteRepository creates a new InMemoryVideoFavoriteRepository.
func NewInMemoryVideoFavoriteRepository() *InMemoryVideoFavoriteRepository {
	return &InMemoryVideoFavoriteRepository{
		favorites: make(map[uint64]map[uint64]struct{}),
	}
}

// Add records a user's like for a video in memory.
func (r *InMemoryVideoFavoriteRepository) Add(ctx context.Context, favorite *model.VideoFavorite) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.favorites[favorite.UserID]; !ok {
		r.favorites[favorite.UserID] = make(map[uint64]struct{})
	}
	r.favorites[favorite.UserID][favorite.VideoID] = struct{}{}
	return nil
}

// Remove deletes a user's like for a video from memory.
func (r *InMemoryVideoFavoriteRepository) Remove(ctx context.Context, favorite *model.VideoFavorite) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if userFavorites, ok := r.favorites[favorite.UserID]; ok {
		delete(userFavorites, favorite.VideoID)
	}
	return nil
}

// ListVideoIDsByUserID returns a list of video IDs a user has favorited from memory.
func (r *InMemoryVideoFavoriteRepository) ListVideoIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userFavorites, ok := r.favorites[userID]
	if !ok {
		return []uint64{}, nil
	}

	videoIDs := make([]uint64, 0, len(userFavorites))
	for videoID := range userFavorites {
		videoIDs = append(videoIDs, videoID)
	}

	return videoIDs, nil
}
