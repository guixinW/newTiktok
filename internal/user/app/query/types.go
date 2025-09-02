package query

import "time"

type User struct {
	UUID           string
	Name           string
	Age            uint16
	Gender         uint16
	FollowingCount uint64
	FollowerCount  uint64
	TotalFavorite  uint64
	WorkCount      uint64
	FavoriteCount  uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
