package persistence

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log/slog"
	"newTiktoken/internal/user/domain/model"
	"newTiktoken/internal/user/domain/repository"
)

// GormUser 是用于 GORM 的用户模型
type GormUser struct {
	ID             uint64 `gorm:"primaryKey"`
	Username       string `gorm:"type:varchar(255);uniqueIndex"`
	HashedPassword string
	FollowingCount uint64
	FollowerCount  uint64
	TotalFavorited uint64
	WorkCount      uint64
	FavoriteCount  uint64
	CreatedAt      int64 `gorm:"autoCreateTime"`
	UpdatedAt      int64 `gorm:"autoUpdateTime"`
}

func (GormUser) TableName() string {
	return "users"
}

// toDomainModel 将 GORM 模型转换为领域模型
func (g *GormUser) toDomainModel() *model.User {
	return &model.User{
		ID:             g.ID,
		Username:       g.Username,
		HashedPassword: g.HashedPassword,
		FollowingCount: g.FollowingCount,
		FollowerCount:  g.FollowerCount,
		TotalFavorited: g.TotalFavorited,
		WorkCount:      g.WorkCount,
		FavoriteCount:  g.FavoriteCount,
	}
}

// fromDomainModel 将领域模型转换为 GORM 模型
func fromDomainModel(u *model.User) *GormUser {
	return &GormUser{
		ID:             u.ID,
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		FollowingCount: u.FollowingCount,
		FollowerCount:  u.FollowerCount,
		TotalFavorited: u.TotalFavorited,
		WorkCount:      u.WorkCount,
		FavoriteCount:  u.FavoriteCount,
	}
}

type mysqlUserRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewMySQLUserRepository 创建一个新的 MySQL 用户仓库实例
func NewMySQLUserRepository(db *gorm.DB, log *slog.Logger) (repository.UserRepository, error) {
	// 自动迁移模式，确保表结构是最新的
	if err := db.AutoMigrate(&GormUser{}); err != nil {
		log.Error("gorm auto migration failed", "table", GormUser{}.TableName(), "error", err)
		return nil, err
	}
	return &mysqlUserRepository{
		db:  db,
		log: log.With("repo", "mysql"),
	}, nil
}

// Save 保存用户 (创建或更新)
func (r *mysqlUserRepository) Save(ctx context.Context, user *model.User) error {
	gormUser := fromDomainModel(user)
	err := r.db.WithContext(ctx).Save(gormUser).Error
	if err != nil {
		r.log.Error("failed to save user", "user_id", user.ID, "username", user.Username, "error", err)
	}
	// 在 GORM 中，Save 会更新 user 的 ID，我们需要把它传回领域模型
	user.ID = gormUser.ID
	return err
}

// FindByUsername 根据用户名查找用户
func (r *mysqlUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var gormUser GormUser
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("user not found by username", "username", username)
			return nil, nil
		}
		r.log.Error("failed to find user by username", "username", username, "error", err)
		return nil, err
	}
	return gormUser.toDomainModel(), nil
}

// FindByID 根据用户ID查找用户
func (r *mysqlUserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var gormUser GormUser
	err := r.db.WithContext(ctx).First(&gormUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("user not found by id", "user_id", id)
			return nil, nil
		}
		r.log.Error("failed to find user by id", "user_id", id, "error", err)
		return nil, err
	}
	return gormUser.toDomainModel(), nil
}
