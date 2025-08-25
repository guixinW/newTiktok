package grpc

import (
	"context"
	"log/slog"
	"newTiktoken/internal/user_relation/application"
	"newTiktoken/pkg/pb/user_relation"
)

// UserRelationServer is the gRPC server for the user relation service.
type UserRelationServer struct {
	user_relation.UnimplementedRelationServiceServer
	appSvc *application.UserRelationApplicationService
	logger *slog.Logger
}

// NewUserRelationServer creates a new UserRelationServer.
func NewUserRelationServer(appSvc *application.UserRelationApplicationService) *UserRelationServer {
	return &UserRelationServer{appSvc: appSvc}
}

// RelationAction implements the gRPC RelationAction method.
func (s *UserRelationServer) RelationAction(ctx context.Context, req *user_relation.RelationActionRequest) (*user_relation.RelationActionResponse, error) {
	// TODO: Implement request handling
	return &user_relation.RelationActionResponse{}, nil
}

func (s *UserRelationServer) RelationFollowerList(ctx context.Context, req *user_relation.RelationFollowerListRequest) (*user_relation.RelationFollowerListResponse, error) {
	// TODO: Implement request handling
	return &user_relation.RelationFollowerListResponse{}, nil
}

func (s *UserRelationServer) RelationFollowList(ctx context.Context, req *user_relation.RelationFollowListRequest) (*user_relation.RelationFollowListResponse, error) {
	// TODO: Implement request handling
	return &user_relation.RelationFollowListResponse{}, nil
}

func (s *UserRelationServer) RelationFriendList(ctx context.Context, req *user_relation.RelationFriendListRequest) (*user_relation.RelationFriendListResponse, error) {
	// TODO: Implement request handling
	return &user_relation.RelationFriendListResponse{}, nil
}
