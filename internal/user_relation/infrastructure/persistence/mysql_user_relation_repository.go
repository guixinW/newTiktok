package persistence

import (
	"context"
	"database/sql"
	"newTiktoken/internal/user_relation/domain/model"
	"newTiktoken/internal/user_relation/domain/repository"
)

// MySQLUserRelationRepository is the MySQL implementation of the UserRelationRepository.
type MySQLUserRelationRepository struct {
	db *sql.DB
}

// NewMySQLUserRelationRepository creates a new MySQLUserRelationRepository.
func NewMySQLUserRelationRepository(db *sql.DB) repository.UserRelationRepository {
	return &MySQLUserRelationRepository{db: db}
}

// Follow creates a new follow relation in the database.
func (r *MySQLUserRelationRepository) Follow(ctx context.Context, userID, followerID uint64) error {
	// TODO: Implement database insertion
	return nil
}

// Unfollow deletes a follow relation from the database.
func (r *MySQLUserRelationRepository) Unfollow(ctx context.Context, userID, followerID uint64) error {
	// TODO: Implement database deletion
	return nil
}

// GetFollowers gets followers of a user from the database.
func (r *MySQLUserRelationRepository) GetFollowers(ctx context.Context, userID uint64) ([]*model.UserRelation, error) {
	// TODO: Implement database query
	return nil, nil
}

// GetFollowing gets following of a user from the database.
func (r *MySQLUserRelationRepository) GetFollowing(ctx context.Context, userID uint64) ([]*model.UserRelation, error) {
	// TODO: Implement database query
	return nil, nil
}
