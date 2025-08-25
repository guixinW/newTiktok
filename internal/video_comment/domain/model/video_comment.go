package model

import "time"

// Comment represents the domain model for a comment.
// It is independent of the database or gRPC layers.
type Comment struct {
	ID        uint64
	UserID    uint64
	VideoID   uint64
	Content   string
	CreatedAt time.Time
}
