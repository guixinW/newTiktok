package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"newTiktoken/internal/user/domain/model"
	"newTiktoken/internal/user/domain/repository"
	"time"
)

const (
	userCacheKeyByID       = "user:id:%d"
	userCacheKeyByUsername = "user:username:%s"
	userCacheDuration      = time.Hour
)

type redisUserRepository struct {
	redisClient *redis.Client
	next        repository.UserRepository // 下一个仓库 (例如 MySQL)
}

// NewRedisUserRepository 创建一个新的 Redis 用户仓库 (缓存层)
func NewRedisUserRepository(redisClient *redis.Client, next repository.UserRepository) repository.UserRepository {
	return &redisUserRepository{
		redisClient: redisClient,
		next:        next,
	}
}

// Save 保存用户，并使缓存失效
func (r *redisUserRepository) Save(ctx context.Context, user *model.User) error {
	// 先调用下一个仓库保存到持久化存储
	if err := r.next.Save(ctx, user); err != nil {
		return err
	}

	// 持久化成功后，使相关缓存失效
	// 我们需要确保 user.ID 是有效的
	if user.ID > 0 {
		keyID := fmt.Sprintf(userCacheKeyByID, user.ID)
		r.redisClient.Del(ctx, keyID)
	}
	if user.Username != "" {
		keyUsername := fmt.Sprintf(userCacheKeyByUsername, user.Username)
		r.redisClient.Del(ctx, keyUsername)
	}

	return nil
}

// FindByUsername 首先尝试从 Redis 缓存中查找用户，如果未找到，则从下一个仓库中查找
func (r *redisUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	key := fmt.Sprintf(userCacheKeyByUsername, username)
	return r.findUserFromCacheOrNext(ctx, key, func(ctx context.Context) (*model.User, error) {
		return r.next.FindByUsername(ctx, username)
	})
}

// FindByID 首先尝试从 Redis 缓存中查找用户，如果未找到，则从下一个仓库中查找
func (r *redisUserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	key := fmt.Sprintf(userCacheKeyByID, id)
	return r.findUserFromCacheOrNext(ctx, key, func(ctx context.Context) (*model.User, error) {
		return r.next.FindByID(ctx, id)
	})
}

// findUserFromCacheOrNext 是一个通用函数，用于处理缓存逻辑
func (r *redisUserRepository) findUserFromCacheOrNext(ctx context.Context, key string, findNext func(context.Context) (*model.User, error)) (*model.User, error) {
	// 1. 尝试从 Redis 获取
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == nil {
		// 缓存命中
		var user model.User
		if json.Unmarshal([]byte(val), &user) == nil {
			return &user, nil
		}
	}

	if !errors.Is(redis.Nil, err) {
		// 如果是除了 "not found" 之外的其他 Redis 错误，则记录并继续
		// log.Printf("Redis error on get: %v", err)
	}

	// 2. 缓存未命中，从下一个仓库 (MySQL) 获取
	user, err := findNext(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// 记录未找到，但为了防止缓存穿透，可以缓存一个空的标记
		// 这里为了简化，我们不缓存空结果
		return nil, nil
	}

	// 3. 将从数据库获得的数据序列化并存入 Redis
	jsonData, err := json.Marshal(user)
	if err == nil {
		r.redisClient.Set(ctx, key, jsonData, userCacheDuration)
	}

	return user, nil
}
