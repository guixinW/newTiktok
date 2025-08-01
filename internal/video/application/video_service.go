package application

import (
	"context"
	"newTiktoken/internal/video/domain/model"
	"newTiktoken/internal/video/domain/repository"
	"newTiktoken/pkg/userpb"
	"time"
)

// VideoService provides video-related application services.
type VideoService struct {
	videoRepo repository.VideoRepository
	// In a real-world scenario, you might have a user service client
	// to get author info, check favorites, etc.
	// userServiceClient userpb.UserServiceClient
}

// NewVideoService creates a new VideoService.
func NewVideoService(videoRepo repository.VideoRepository) *VideoService {
	return &VideoService{videoRepo: videoRepo}
}

// Feed retrieves the video feed.
func (s *VideoService) Feed(ctx context.Context, latestTime int64, tokenUserID uint64) ([]*model.Video, int64, error) {
	// Business logic for the feed.
	// For now, it's a direct call to the repository.
	// In a real implementation, you might add more complex logic here,
	// like personalizing the feed for the user.
	videos, nextTime, err := s.videoRepo.Feed(ctx, latestTime)
	if err != nil {
		return nil, 0, err
	}

	// Here you would enrich the video data, for example, by checking if the tokenUserID
	// has favorited any of the videos. This is a placeholder for that logic.
	for _, video := range videos {
		// This is just an example, in a real app you would call a favorite service.
		// video.IsFavorite = someFavoriteCheck(tokenUserID, video.ID)
	}

	return videos, nextTime, nil
}

// PublishAction handles the creation of a new video.
func (s *VideoService) PublishAction(ctx context.Context, authorID uint64, playURL, title string) error {
	// In a real app, you would fetch complete author details from a user service.
	author := &userpb.User{Id: authorID} // Placeholder, only ID is known from token.

	video := &model.Video{
		Author:   author,
		PlayURL:  playURL,
		Title:    title,
		CreateAt: uint64(time.Now().Unix()),
		// Other fields like FavoriteCount, CommentCount would be initialized to 0.
	}
	return s.videoRepo.Create(ctx, video)
}

// PublishList retrieves the list of videos published by a user.
func (s *VideoService) PublishList(ctx context.Context, userID uint64, tokenUserID uint64) ([]*model.Video, error) {
	// Business logic for the publish list.
	// You might add logic here to check if the tokenUserID can view the userID's publish list.
	videos, err := s.videoRepo.ListByAuthorID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Enrich video data, e.g., check for favorites by tokenUserID.
	for _, video := range videos {
		// This is just an example, in a real app you would call a favorite service.
		// video.IsFavorite = someFavoriteCheck(tokenUserID, video.ID)
	}
	return videos, nil
}
