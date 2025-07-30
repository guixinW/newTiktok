package persistence

import (
	"context"
	"fmt"
	"newTiktoken/internal/user/domain/model"
	"newTiktoken/internal/user/domain/repository"
	"sync"
)

// ErrUserNotFound 是一个标准错误，表示未找到用户
var ErrUserNotFound = fmt.Errorf("user not found")

// InMemoryUserRepository 是 UserRepository 的一个内存实现，用于测试和开发
type InMemoryUserRepository struct {
	mu      sync.RWMutex
	users   map[uint64]*model.User
	nextID  uint64
}

// NewInMemoryUserRepository 创建一个新的内存用户仓库
func NewInMemoryUserRepository() repository.UserRepository {
	return &InMemoryUserRepository{
		users:  make(map[uint64]*model.User),
		nextID: 1,
	}
}

// Save 保存用户。如果用户ID为0，则认为是新用户并分配ID。
func (r *InMemoryUserRepository) Save(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID == 0 {
		user.ID = r.nextID
		r.nextID++
	}
	r.users[user.ID] = user
	return nil
}

// FindByUsername 通过用户名查找用户
func (r *InMemoryUserRepository) FindByUsername(_ context.Context, username string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// FindByID 通过ID查找用户
func (r *InMemoryUserRepository) FindByID(_ context.Context, id uint64) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}
