package model

import "time"

// UserRelation represents the domain model for a user relation.
// It is independent of the database or gRPC layers.
type UserRelation struct {
	ID         uint64
	UserID     uint64
	FollowerID uint64
	CreatedAt  time.Time
}
