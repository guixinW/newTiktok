package user

import (
	"context"
)

// Repository 是user domain repository的接口
type Repository interface {
	GetUser(ctx context.Context, userUUID string) (*User, error)
	AddUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, userUUID string, updateFn func(
		ctx context.Context,
		user *User,
	) (*User, error)) error
}
