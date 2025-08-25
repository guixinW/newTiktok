package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"newTiktoken/internal/video_favorite/domain/model"
	"newTiktoken/internal/video_favorite/domain/repository"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisVideoFavoriteRepository is a cached repository for video favorites.
type RedisVideoFavoriteRepository struct {
	client *redis.Client
	next   repository.VideoFavoriteRepository
	ttl    time.Duration
}

// NewRedisVideoFavoriteRepository creates a new RedisVideoFavoriteRepository.
func NewRedisVideoFavoriteRepository(client *redis.Client, next repository.VideoFavoriteRepository) *RedisVideoFavoriteRepository {
	return &RedisVideoFavoriteRepository{
		client: client,
		next:   next,
		ttl:    1 * time.Hour,
	}
}

// Helper function to generate a cache key for a user's favorite list.
func userFavoriteListKey(userID uint64) string {
	return fmt.Sprintf("video_favorites:user:%d", userID)
}

// Add calls the next repository and then invalidates the cache.
func (r *RedisVideoFavoriteRepository) Add(ctx context.Context, favorite *model.VideoFavorite) error {
	err := r.next.Add(ctx, favorite)
	if err != nil {
		return err
	}

	cacheKey := userFavoriteListKey(favorite.UserID)
	r.client.Del(ctx, cacheKey) // Invalidate cache

	return nil
}

// Remove calls the next repository and then invalidates the cache.
func (r *RedisVideoFavoriteRepository) Remove(ctx context.Context, favorite *model.VideoFavorite) error {
	err := r.next.Remove(ctx, favorite)
	if err != nil {
		return err
	}

	cacheKey := userFavoriteListKey(favorite.UserID)
	r.client.Del(ctx, cacheKey) // Invalidate cache

	return nil
}

// ListVideoIDsByUserID tries to get the list from Redis first.
func (r *RedisVideoFavoriteRepository) ListVideoIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error) {
	cacheKey := userFavoriteListKey(userID)

	// 1. Try to get from cache
	result, err := r.client.Get(ctx, cacheKey).Result()
	if err == nil {
		var videoIDs []uint64
		if json.Unmarshal([]byte(result), &videoIDs) == nil {
			return videoIDs, nil
		}
	}

	// 2. Cache miss, get from the next repository
	videoIDs, err := r.next.ListVideoIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Cache the result
	if len(videoIDs) > 0 {
		data, err := json.Marshal(videoIDs)
		if err == nil {
			r.client.Set(ctx, cacheKey, data, r.ttl)
		}
	}

	return videoIDs, nil
}
