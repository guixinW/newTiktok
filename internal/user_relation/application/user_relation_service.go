package application

import (
	"context"
	"newTiktoken/internal/user_relation/domain/repository"
	"newTiktoken/pkg/pb/user_relation"
)

// UserRelationApplicationService provides user relation-related application services.
type UserRelationApplicationService struct {
	userRelationRepo repository.UserRelationRepository
}

// NewUserRelationApplicationService creates a new UserRelationApplicationService.
func NewUserRelationApplicationService(userRelationRepo repository.UserRelationRepository) *UserRelationApplicationService {
	return &UserRelationApplicationService{userRelationRepo: userRelationRepo}
}

// Follow handles the follow action.
func (s *UserRelationApplicationService) Follow(ctx context.Context, req *user_relation.FollowRequest) error {
	// TODO: Implement business logic
	return nil
}

// Unfollow handles the unfollow action.
func (s *UserRelationApplicationService) Unfollow(ctx context.Context, req *user_relation.UnfollowRequest) error {
	// TODO: Implement business logic
	return nil
}

// GetFollowers handles the get followers query.
func (s *UserRelationApplicationService) GetFollowers(ctx context.Context, req *user_relation.GetFollowersRequest) (*user_relation.GetFollowersResponse, error) {
	// TODO: Implement business logic
	return nil, nil
}

// GetFollowing handles the get following query.
func (s *UserRelationApplicationService) GetFollowing(ctx context.Context, req *user_relation.GetFollowingRequest) (*user_relation.GetFollowingResponse, error) {
	// TODO: Implement business logic
	return nil, nil
}
