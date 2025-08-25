package grpc

import (
	"context"
	"log/slog"
	"newTiktoken/internal/video_comment/application"
	"newTiktoken/pkg/pb/video_comment"
)

// CommentServer is the gRPC server for the comment service.
type CommentServer struct {
	video_comment.UnimplementedCommentServiceServer
	appSvc *application.CommentApplicationService
	logger *slog.Logger
}

// NewCommentServer creates a new CommentServer.
func NewCommentServer(appSvc *application.CommentApplicationService, logger *slog.Logger) *CommentServer {
	return &CommentServer{
		appSvc: appSvc,
		logger: logger,
	}
}

func (s *CommentServer) CommentAction(ctx context.Context, req *video_comment.CommentActionRequest) (*video_comment.CommentActionResponse, error) {
	// TODO: Implement request handling
	return &video_comment.CommentActionResponse{}, nil
}

func (s *CommentServer) CommentList(ctx context.Context, req *video_comment.CommentListRequest) (*video_comment.CommentListResponse, error) {
	// TODO: Implement request handling
	return &video_comment.CommentListResponse{}, nil
}
