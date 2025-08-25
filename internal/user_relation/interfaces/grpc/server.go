package grpc

import (
	"context"
	"newTiktoken/internal/user_relation/application"
	"newTiktoken/pkg/pb/user_relation"
)

// UserRelationServer is the gRPC server for the user relation service.
type UserRelationServer struct {
	user_relation.UnimplementedUserRelationServiceServer
	appSvc *application.UserRelationApplicationService
}

// NewUserRelationServer creates a new UserRelationServer.
func NewUserRelationServer(appSvc *application.UserRelationApplicationService) *UserRelationServer {
	return &UserRelationServer{appSvc: appSvc}
}

// Follow implements the gRPC Follow method.
func (s *UserRelationServer) Follow(ctx context.Context, req *user_relation.FollowRequest) (*user_relation.FollowResponse, error) {
	// TODO: Implement request handling
	return &user_relation.FollowResponse{},
		nil,
}

// Unfollow implements the gRPC Unfollow method.
func (s *UserRelationServer) Unfollow(ctx context.Context, req *user_relation.UnfollowRequest) (*user_relation.UnfollowResponse, error) {
	// TODO: Implement request handling
	return &user_relation.UnfollowResponse{},
		nil,
}

// GetFollowers implements the gRPC GetFollowers method.
func (s *UserRelationServer) GetFollowers(ctx context.Context, req *user_relation.GetFollowersRequest) (*user_relation.GetFollowersResponse, error) {
	// TODO: Implement request handling
	return &user_relation.GetFollowersResponse{},
		nil,
}

// GetFollowing implements the gRPC GetFollowing method.
func (s *UserRelationServer) GetFollowing(ctx context.Context, req *user_relation.GetFollowingRequest) (*user_relation.GetFollowingResponse, error) {
	// TODO: Implement request handling
	return &user_relation.GetFollowingResponse{},
		nil,
}
