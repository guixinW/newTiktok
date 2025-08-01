package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"newTiktoken/internal/video/domain/model"
	"newTiktoken/internal/video/domain/repository"
)

// RedisVideoRepository is a cached repository that wraps another VideoRepository (e.g., MySQL).
// It uses Redis to cache video data.
type RedisVideoRepository struct {
	client *redis.Client
	next   repository.VideoRepository // The next repository in the chain (e.g., MySQL)
	ttl    time.Duration              // Time-to-live for cache entries
}

// NewRedisVideoRepository creates a new RedisVideoRepository.
func NewRedisVideoRepository(client *redis.Client, next repository.VideoRepository) *RedisVideoRepository {
	return &RedisVideoRepository{
		client: client,
		next:   next,
		ttl:    1 * time.Hour, // Default cache TTL of 1 hour
	}
}

// Helper function to generate a cache key for a user's video list.
func userVideoListKey(authorID uint64) string {
	return fmt.Sprintf("videos:user:%d", authorID)
}

// Feed attempts to get the feed from the cache, otherwise falls back to the next repository.
// Caching the main feed is complex and often requires a more sophisticated strategy
// (e.g., using Redis sorted sets). For simplicity, this implementation will bypass caching for the feed
// and go directly to the source. A real-world implementation would be more involved.
func (r *RedisVideoRepository) Feed(ctx context.Context, latestTime int64) ([]*model.Video, int64, error) {
	// Bypassing cache for the main feed for simplicity.
	// A proper implementation would likely involve a Redis Sorted Set keyed by timestamp.
	return r.next.Feed(ctx, latestTime)
}

// Create clears the cache for the author's video list and then calls the next repository to create the video.
func (r *RedisVideoRepository) Create(ctx context.Context, video *model.Video) error {
	// First, perform the operation in the persistent store.
	err := r.next.Create(ctx, video)
	if err != nil {
		return err
	}

	// If successful, invalidate the cache for this user's video list.
	cacheKey := userVideoListKey(video.Author.Id)
	r.client.Del(ctx, cacheKey) // Fire and forget deletion

	return nil
}

// ListByAuthorID first tries to get the video list from Redis. If not found,
// it gets it from the next repository and caches the result.
func (r *RedisVideoRepository) ListByAuthorID(ctx context.Context, authorID uint64) ([]*model.Video, error) {
	cacheKey := userVideoListKey(authorID)

	// 1. Try to get from cache
	result, err := r.client.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		var videos []*model.Video
		if json.Unmarshal([]byte(result), &videos) == nil {
			return videos, nil
		}
	}

	// 2. Cache miss or error, get from the next repository (e.g., MySQL)
	videos, err := r.next.ListByAuthorID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	// 3. Cache the result for future requests
	if len(videos) > 0 {
		data, err := json.Marshal(videos)
		if err == nil {
			r.client.Set(ctx, cacheKey, data, r.ttl)
		}
	}

	return videos, nil
}
