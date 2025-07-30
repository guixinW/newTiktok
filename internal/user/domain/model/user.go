package model

import (
	"time"
)

// User 是领域模型，代表一个用户
type User struct {
	ID             uint64
	Username       string
	HashedPassword string
	FollowingCount uint64
	FollowerCount  uint64
	TotalFavorited uint64
	WorkCount      uint64
	FavoriteCount  uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewUser 创建一个新的用户实例
func NewUser(username, hashedPassword string) *User {
	return &User{
		Username:       username,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CheckPassword 验证密码是否正确 (伪代码，实际应使用 bcrypt)
func (u *User) CheckPassword(password string) bool {
	// 在实际应用中，这里应该使用 bcrypt.CompareHashAndPassword
	// 为了简化，我们这里只做简单比较
	// return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)) == nil
	return u.HashedPassword == password // 仅为示例
}
