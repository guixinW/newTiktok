package persistence

import (
	"context"
	"errors"
	"log/slog"
	"newTiktoken/internal/video/domain/model"
	"newTiktoken/internal/video/domain/repository"
	"newTiktoken/pkg/pb/user"
	"time"

	"gorm.io/gorm"
)

// GormUser is a slimmed-down copy from the user persistence package.
// In a real-world scenario, this might be in a shared package.
type GormUser struct {
	ID             uint64 `gorm:"primaryKey"`
	Username       string `gorm:"type:varchar(255);uniqueIndex"`
	FollowingCount uint64
	FollowerCount  uint64
	TotalFavorited uint64
	WorkCount      uint64
	FavoriteCount  uint64
}

func (GormUser) TableName() string {
	return "users"
}

// GormVideo is the GORM model for a video.
type GormVideo struct {
	ID            uint64   `gorm:"primaryKey"`
	AuthorID      uint64   `gorm:"index"`
	Author        GormUser `gorm:"foreignKey:AuthorID"`
	PlayURL       string   `gorm:"type:varchar(255)"`
	FavoriteCount uint64
	CommentCount  uint64
	Title         string `gorm:"type:varchar(255)"`
	ShareCount    uint64
	CreatedAt     int64 `gorm:"autoCreateTime"`
}

func (GormVideo) TableName() string {
	return "videos"
}

// toDomainModel converts the GORM video model to a domain model.
// It also converts the nested GormUser to a userpb.User.
func (g *GormVideo) toDomainModel() *model.Video {
	return &model.Video{
		ID: g.ID,
		Author: &user.User{
			Id:             g.Author.ID,
			Name:           g.Author.Username, // Correctly map Username to Name
			FollowingCount: g.Author.FollowingCount,
			FollowerCount:  g.Author.FollowerCount,
			TotalFavorite:  g.Author.TotalFavorited, // Correctly map TotalFavorited to TotalFavorite
			WorkCount:      g.Author.WorkCount,
			FavoriteCount:  g.Author.FavoriteCount,
		},
		PlayURL:       g.PlayURL,
		FavoriteCount: g.FavoriteCount,
		CommentCount:  g.CommentCount,
		Title:         g.Title,
		ShareCount:    g.ShareCount,
		CreateAt:      uint64(g.CreatedAt),
	}
}

// fromDomainModel converts a domain video model to a GORM model.
func fromDomainModel(v *model.Video) *GormVideo {
	return &GormVideo{
		ID:            v.ID,
		AuthorID:      v.Author.Id,
		PlayURL:       v.PlayURL,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		Title:         v.Title,
		ShareCount:    v.ShareCount,
		CreatedAt:     int64(v.CreateAt),
	}
}

// MySQLVideoRepository is a GORM-based implementation of the VideoRepository.
type MySQLVideoRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewMySQLVideoRepository creates a new MySQLVideoRepository.
func NewMySQLVideoRepository(db *gorm.DB, log *slog.Logger) (repository.VideoRepository, error) {
	if err := db.AutoMigrate(&GormVideo{}); err != nil {
		log.Error("gorm auto migration failed", "table", GormVideo{}.TableName(), "error", err)
		return nil, err
	}
	return &MySQLVideoRepository{
		db:  db,
		log: log.With("repo", "mysql_video"),
	}, nil
}

// Feed retrieves a list of videos from the database, ordered by creation time.
func (r *MySQLVideoRepository) Feed(ctx context.Context, latestTime int64) ([]*model.Video, int64, error) {
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}

	var gormVideos []GormVideo
	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("created_at < ?", latestTime).
		Order("created_at DESC").
		Limit(30).
		Find(&gormVideos).Error

	if err != nil {
		r.log.Error("failed to query feed", "error", err)
		return nil, 0, err
	}

	videos := make([]*model.Video, len(gormVideos))
	for i, gv := range gormVideos {
		videos[i] = gv.toDomainModel()
	}

	var newNextTime int64
	if len(videos) > 0 {
		newNextTime = int64(videos[len(videos)-1].CreateAt)
	}

	return videos, newNextTime, nil
}

// Create inserts a new video record into the database.
func (r *MySQLVideoRepository) Create(ctx context.Context, video *model.Video) error {
	gormVideo := fromDomainModel(video)
	err := r.db.WithContext(ctx).Create(gormVideo).Error
	if err != nil {
		r.log.Error("failed to create video", "author_id", video.Author.Id, "title", video.Title, "error", err)
		return err
	}
	video.ID = gormVideo.ID // Update domain model with new ID
	return nil
}

// ListByAuthorID retrieves all videos for a specific author from the database.
func (r *MySQLVideoRepository) ListByAuthorID(ctx context.Context, authorID uint64) ([]*model.Video, error) {
	var gormVideos []GormVideo
	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("author_id = ?", authorID).
		Order("created_at DESC").
		Find(&gormVideos).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("no videos found for author", "author_id", authorID)
			return []*model.Video{}, nil // Return empty slice instead of nil
		}
		r.log.Error("failed to list videos by author", "author_id", authorID, "error", err)
		return nil, err
	}

	videos := make([]*model.Video, len(gormVideos))
	for i, gv := range gormVideos {
		videos[i] = gv.toDomainModel()
	}

	return videos, nil
}
