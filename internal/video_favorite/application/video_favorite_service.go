package application

import (
	"context"
	"errors"
	videoRepo "newTiktoken/internal/video/domain/repository"
	"newTiktoken/internal/video_favorite/domain/model"
	"newTiktoken/internal/video_favorite/domain/repository"
	"newTiktoken/pkg/pb/video"
	"newTiktoken/pkg/pb/video_favorite"
)

// VideoFavoriteService provides favorite-related application services.
type VideoFavoriteService struct {
	favoriteRepo repository.VideoFavoriteRepository
	videoRepo    videoRepo.VideoRepository // To fetch video details
}

// NewVideoFavoriteService creates a new VideoFavoriteService.
func NewVideoFavoriteService(favoriteRepo repository.VideoFavoriteRepository, videoRepo videoRepo.VideoRepository) *VideoFavoriteService {
	return &VideoFavoriteService{
		favoriteRepo: favoriteRepo,
		videoRepo:    videoRepo,
	}
}

// FavoriteAction handles liking and unliking videos.
func (s *VideoFavoriteService) FavoriteAction(ctx context.Context, userID uint64, videoID uint64, actionType video_favorite.VideoActionType) error {
	fav := &model.VideoFavorite{
		UserID:  userID,
		VideoID: videoID,
	}

	switch actionType {
	case video_favorite.VideoActionType_LIKE:
		return s.favoriteRepo.Add(ctx, fav)
	case video_favorite.VideoActionType_CANCEL_LIKE:
		return s.favoriteRepo.Remove(ctx, fav)
	case video_favorite.VideoActionType_DISLIKE, video_favorite.VideoActionType_CANCEL_DISLIKE:
		return errors.New("dislike functionality not implemented")
	default:
		return errors.New("invalid favorite action type")
	}
}

// FavoriteList retrieves the list of videos favorited by a user.
func (s *VideoFavoriteService) FavoriteList(ctx context.Context, userID uint64, tokenUserID uint64) ([]*video.Video, error) {
	videoIDs, err := s.favoriteRepo.ListVideoIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(videoIDs) == 0 {
		return []*video.Video{}, nil
	}

	// This is a placeholder. In a real microservices architecture, you would make a gRPC call
	// to the video service to get details for these video IDs.
	videos := make([]*video.Video, len(videoIDs))
	for i, id := range videoIDs {
		videos[i] = &video.Video{Id: id, Title: "Faked Video Title"} // Fake data
	}

	return videos, nil
}
