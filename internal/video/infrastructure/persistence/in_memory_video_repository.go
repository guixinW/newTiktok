package persistence

import (
	"context"
	"sort"
	"sync"
	"time"

	"newTiktoken/internal/video/domain/model"
)

// InMemoryVideoRepository is an in-memory implementation of VideoRepository for testing and development.
// It holds data in memory and is not suitable for production.
type InMemoryVideoRepository struct {
	mu     sync.RWMutex
	videos map[uint64]*model.Video
	nextID uint64
}

// NewInMemoryVideoRepository creates a new InMemoryVideoRepository.
func NewInMemoryVideoRepository() *InMemoryVideoRepository {
	return &InMemoryVideoRepository{
		videos: make(map[uint64]*model.Video),
		nextID: 1,
	}
}

// Feed returns a list of videos from memory.
func (r *InMemoryVideoRepository) Feed(ctx context.Context, latestTime int64) ([]*model.Video, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}

	var videoList []*model.Video
	for _, v := range r.videos {
		if v.CreateAt < uint64(latestTime) {
			// In a real implementation, you'd create a deep copy
			videoList = append(videoList, v)
		}
	}

	// Sort by creation time descending to get the latest videos first
	sort.Slice(videoList, func(i, j int) bool {
		return videoList[i].CreateAt > videoList[j].CreateAt
	})

	// Limit to a fixed number, e.g., 30, as is common in feeds
	if len(videoList) > 30 {
		videoList = videoList[:30]
	}

	var newNextTime int64
	if len(videoList) > 0 {
		// The next time for the client to query is the timestamp of the oldest video in this batch
		newNextTime = int64(videoList[len(videoList)-1].CreateAt)
	}

	return videoList, newNextTime, nil
}

// Create saves a new video to memory.
func (r *InMemoryVideoRepository) Create(ctx context.Context, video *model.Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a copy to store in the map to avoid external modifications
	videoToStore := *video
	videoToStore.ID = r.nextID
	r.videos[r.nextID] = &videoToStore
	r.nextID++
	return nil
}

// ListByAuthorID returns a list of videos for a given author from memory.
func (r *InMemoryVideoRepository) ListByAuthorID(ctx context.Context, authorID uint64) ([]*model.Video, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var videoList []*model.Video
	for _, v := range r.videos {
		if v.Author != nil && v.Author.Id == authorID {
			videoList = append(videoList, v)
		}
	}

	// Sort by creation time descending
	sort.Slice(videoList, func(i, j int) bool {
		return videoList[i].CreateAt > videoList[j].CreateAt
	})

	return videoList, nil
}
