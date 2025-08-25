package persistence

import (
	"context"
	"log/slog"
	"newTiktoken/internal/video_favorite/domain/model"
	"newTiktoken/internal/video_favorite/domain/repository"

	"gorm.io/gorm"
)

// GormVideoFavorite is the GORM model for a user-video favorite relationship.
type GormVideoFavorite struct {
	UserID  uint64 `gorm:"primaryKey"`
	VideoID uint64 `gorm:"primaryKey"`
}

func (GormVideoFavorite) TableName() string {
	return "video_favorites"
}

// MySQLVideoFavoriteRepository is a GORM-based implementation of the VideoFavoriteRepository.
type MySQLVideoFavoriteRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewMySQLVideoFavoriteRepository creates a new MySQLVideoFavoriteRepository.
func NewMySQLVideoFavoriteRepository(db *gorm.DB, log *slog.Logger) (repository.VideoFavoriteRepository, error) {
	if err := db.AutoMigrate(&GormVideoFavorite{}); err != nil {
		log.Error("gorm auto migration failed", "table", GormVideoFavorite{}.TableName(), "error", err)
		return nil, err
	}
	return &MySQLVideoFavoriteRepository{
		db:  db,
		log: log.With("repo", "mysql_video_favorite"),
	}, nil
}

// Add inserts a new favorite record into the database.
func (r *MySQLVideoFavoriteRepository) Add(ctx context.Context, favorite *model.VideoFavorite) error {
	gormFavorite := &GormVideoFavorite{UserID: favorite.UserID, VideoID: favorite.VideoID}
	err := r.db.WithContext(ctx).Create(gormFavorite).Error
	if err != nil {
		r.log.Error("failed to add video favorite", "user_id", favorite.UserID, "video_id", favorite.VideoID, "error", err)
		return err
	}
	return nil
}

// Remove deletes a favorite record from the database.
func (r *MySQLVideoFavoriteRepository) Remove(ctx context.Context, favorite *model.VideoFavorite) error {
	gormFavorite := &GormVideoFavorite{UserID: favorite.UserID, VideoID: favorite.VideoID}
	result := r.db.WithContext(ctx).Delete(gormFavorite)
	if result.Error != nil {
		r.log.Error("failed to remove video favorite", "user_id", favorite.UserID, "video_id", favorite.VideoID, "error", result.Error)
		return result.Error
	}
	return nil
}

// ListVideoIDsByUserID retrieves all video IDs for a specific user from the database.
func (r *MySQLVideoFavoriteRepository) ListVideoIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error) {
	var videoIDs []uint64
	err := r.db.WithContext(ctx).Model(&GormVideoFavorite{}).Where("user_id = ?", userID).Pluck("video_id", &videoIDs).Error
	if err != nil {
		r.log.Error("failed to list favorited video ids by user", "user_id", userID, "error", err)
		return nil, err
	}
	return videoIDs, nil
}
