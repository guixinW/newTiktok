package model

// VideoFavorite represents the relationship of a user favoriting a video.
// It is the core domain model for a "like".
type VideoFavorite struct {
	UserID  uint64
	VideoID uint64
}