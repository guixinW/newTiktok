package domain

import (
	"github.com/pkg/errors"
	"strings"
	"time"
)

type RelationActionType int

const (
	Follow RelationActionType = iota
	Unfollow
	Block
)

// UserRelation 是用户关系领域的核心实体
type UserRelation struct {
	ActivePartyUUID  string
	PassivePartyUUID string
	Status           RelationActionType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewUserRelation 是创建新关注关系的工厂函数
func NewUserRelation(followerID string, followedID string, actionType RelationActionType) (*UserRelation, error) {
	if followerID == "" {
		return nil, errors.New("followerID不能为空")
	}
	if followedID == "" {
		return nil, errors.New("followedID不能为空")
	}
	if strings.Compare(followerID, followedID) == 0 {
		return nil, errors.New("关注者和被关注者不能是同一个人")
	}

	now := time.Now()
	return &UserRelation{
		ActivePartyUUID:  followerID,
		PassivePartyUUID: followedID,
		Status:           actionType,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func (r *UserRelation) Follow() {
	r.Status = Follow
	r.UpdatedAt = time.Now()
}

func (r *UserRelation) Unfollow() {
	r.Status = Unfollow
	r.UpdatedAt = time.Now()
}

func (r *UserRelation) Block() {
	r.Status = Block
	r.UpdatedAt = time.Now()
}
