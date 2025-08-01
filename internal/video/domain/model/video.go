package model

import "newTiktoken/pkg/userpb"

// Video corresponds to the Video message in video.proto.
// This is the core domain model for a video.
type Video struct {
	ID            uint64
	Author        *userpb.User
	PlayURL       string
	FavoriteCount uint64
	CommentCount  uint64
	IsFavorite    bool
	Title         string
	ShareCount    uint64
	CreateAt      uint64
}
