package repository

import (
	"context"
	"newTiktoken/internal/user/domain/model"
)

// UserRepository 是用户仓库的接口定义
// 它定义了用户数据持久化所需的方法
type UserRepository interface {
	// Save 保存一个用户 (创建或更新)
	Save(ctx context.Context, user *model.User) error
	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	// FindByID 根据用户ID查找用户
	FindByID(ctx context.Context, id uint64) (*model.User, error)
}
