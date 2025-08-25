package repository

import (
	"context"
	"newTiktoken/internal/user_relation/domain/model"
)

// UserRelationRepository defines the interface for user relation data operations.
type UserRelationRepository interface {
	Follow(ctx context.Context, userID, followerID uint64) error
	Unfollow(ctx context.Context, userID, followerID uint64) error
	GetFollowers(ctx context.Context, userID uint64) ([]*model.UserRelation, error)
	GetFollowing(ctx context.Context, userID uint64) ([]*model.UserRelation, error)
}
