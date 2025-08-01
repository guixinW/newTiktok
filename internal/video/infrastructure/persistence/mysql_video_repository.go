package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"newTiktoken/internal/video/domain/model"
	"newTiktoken/pkg/userpb"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLVideoRepository is a MySQL-based implementation of the VideoRepository.
type MySQLVideoRepository struct {
	db *sql.DB
}

// NewMySQLVideoRepository creates a new MySQLVideoRepository.
// The dsn (Data Source Name) should be in the format: "user:password@tcp(127.0.0.1:3306)/dbname"
func NewMySQLVideoRepository(dsn string) (*MySQLVideoRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	// You should run a schema migration script here to create the 'videos' table.
	return &MySQLVideoRepository{db: db}, nil
}

// Feed retrieves a list of videos from the database, ordered by creation time.
func (r *MySQLVideoRepository) Feed(ctx context.Context, latestTime int64) ([]*model.Video, int64, error) {
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}

	// Query to get the latest 30 videos before the given timestamp
	query := `
		SELECT id, author_id, play_url, cover_url, title, create_at 
		FROM videos 
		WHERE create_at < FROM_UNIXTIME(?)
		ORDER BY create_at DESC 
		LIMIT 30`

	rows, err := r.db.QueryContext(ctx, query, latestTime)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query feed: %w", err)
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		var v model.Video
		var authorID uint64
		var createAt time.Time
		// Note: cover_url is in the proto but not in the domain model.
		// This is a placeholder for where you would scan it if it were.
		var coverURL string 

		if err := rows.Scan(&v.ID, &authorID, &v.PlayURL, &coverURL, &v.Title, &createAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan video row: %w", err)
		}
		v.Author = &userpb.User{Id: authorID} // In a real app, you'd fetch full author details
		v.CreateAt = uint64(createAt.Unix())
		videos = append(videos, &v)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during video rows iteration: %w", err)
	}

	var newNextTime int64
	if len(videos) > 0 {
		newNextTime = int64(videos[len(videos)-1].CreateAt)
	}

	return videos, newNextTime, nil
}

// Create inserts a new video record into the database.
func (r *MySQLVideoRepository) Create(ctx context.Context, video *model.Video) error {
	query := "INSERT INTO videos (author_id, play_url, title, create_at) VALUES (?, ?, ?, FROM_UNIXTIME(?))"
	_, err := r.db.ExecContext(ctx, query, video.Author.Id, video.PlayURL, video.Title, video.CreateAt)
	if err != nil {
		return fmt.Errorf("failed to insert video: %w", err)
	}
	return nil
}

// ListByAuthorID retrieves all videos for a specific author from the database.
func (r *MySQLVideoRepository) ListByAuthorID(ctx context.Context, authorID uint64) ([]*model.Video, error) {
	query := `
		SELECT id, author_id, play_url, cover_url, title, create_at 
		FROM videos 
		WHERE author_id = ? 
		ORDER BY create_at DESC`

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos by author: %w", err)
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		var v model.Video
		var createAt time.Time
		var coverURL string // Placeholder

		if err := rows.Scan(&v.ID, &v.Author.Id, &v.PlayURL, &coverURL, &v.Title, &createAt); err != nil {
			return nil, fmt.Errorf("failed to scan video row: %w", err)
		}
		v.CreateAt = uint64(createAt.Unix())
		videos = append(videos, &v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during video rows iteration: %w", err)
	}

	return videos, nil
}

// Close closes the database connection.
func (r *MySQLVideoRepository) Close() {
	r.db.Close()
}
