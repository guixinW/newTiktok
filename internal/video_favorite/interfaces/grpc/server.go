package grpc

import (
	"context"
	"newTiktoken/internal/video_favorite/application"
	"newTiktoken/pkg/pb/video_favorite"
)

// Server is the gRPC server for the favorite service.
type Server struct {
	video_favorite.UnimplementedFavoriteServiceServer
	app *application.VideoFavoriteService
}

// NewServer creates a new gRPC server that wraps the application service.
func NewServer(app *application.VideoFavoriteService) *Server {
	return &Server{app: app}
}

// FavoriteAction implements the FavoriteAction RPC endpoint.
func (s *Server) FavoriteAction(ctx context.Context, req *video_favorite.FavoriteActionRequest) (*video_favorite.FavoriteActionResponse, error) {
	err := s.app.FavoriteAction(ctx, req.TokenUserId, req.VideoId, req.ActionType)
	if err != nil {
		return &video_favorite.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, nil
	}

	return &video_favorite.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
	}, nil
}

// FavoriteList implements the FavoriteList RPC endpoint.
func (s *Server) FavoriteList(ctx context.Context, req *video_favorite.FavoriteListRequest) (*video_favorite.FavoriteListResponse, error) {
	videos, err := s.app.FavoriteList(ctx, req.UserId, req.TokenUserId)
	if err != nil {
		return &video_favorite.FavoriteListResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, nil
	}

	return &video_favorite.FavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		VideoList:  videos,
	}, nil
}
