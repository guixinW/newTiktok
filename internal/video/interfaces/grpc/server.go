package grpc

import (
	"context"
	"newTiktoken/internal/video/application"
	"newTiktoken/internal/video/domain/model"
	videopb "newTiktoken/pkg/pb/video"
)

// Server is the gRPC server for the video service. It implements the generated interface.
type Server struct {
	videopb.UnimplementedVideoServiceServer
	app *application.VideoService
}

// NewServer creates a new gRPC server that wraps the application service.
func NewServer(app *application.VideoService) *Server {
	return &Server{app: app}
}

// Feed implements the Feed RPC endpoint.
func (s *Server) Feed(ctx context.Context, req *videopb.FeedRequest) (*videopb.FeedResponse, error) {
	videos, nextTime, err := s.app.Feed(ctx, req.LatestTime, req.TokenUserId)
	if err != nil {
		return &videopb.FeedResponse{
			StatusCode: -1, // Using -1 to indicate an error
			StatusMsg:  err.Error(),
		}, nil // Returning nil error to gRPC framework, as we've encoded it in the response
	}

	return &videopb.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		VideoList:  toProtobufVideos(videos),
		NextTime:   nextTime,
	}, nil
}

// PublishAction implements the PublishAction RPC endpoint.
func (s *Server) PublishAction(ctx context.Context, req *videopb.PublishActionRequest) (*videopb.PublishActionResponse, error) {
	err := s.app.PublishAction(ctx, req.TokenUserId, req.PlayUrl, req.Title)
	if err != nil {
		return &videopb.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, nil
	}

	return &videopb.PublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "Publish successful",
	}, nil
}

// PublishList implements the PublishList RPC endpoint.
func (s *Server) PublishList(ctx context.Context, req *videopb.PublishListRequest) (*videopb.PublishListResponse, error) {
	videos, err := s.app.PublishList(ctx, req.UserId, req.TokenUserId)
	if err != nil {
		return &videopb.PublishListResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, nil
	}

	return &videopb.PublishListResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		VideoList:  toProtobufVideos(videos),
	}, nil
}

// toProtobufVideos is a helper function to convert domain model videos to protobuf videos.
func toProtobufVideos(videos []*model.Video) []*videopb.Video {
	protoVideos := make([]*videopb.Video, len(videos))
	for i, v := range videos {
		protoVideos[i] = &videopb.Video{
			Id:            v.ID,
			Author:        v.Author,
			PlayUrl:       v.PlayURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
			ShareCount:    v.ShareCount,
			CreateAt:      v.CreateAt,
		}
	}
	return protoVideos
}
