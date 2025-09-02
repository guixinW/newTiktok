package user

import (
	"github.com/pkg/errors"
	"time"
)

// User 是领域模型，代表一个用户
type User struct {
	uuid      string
	name      string
	age       uint16
	gender    uint16
	createdAt time.Time
	updatedAt time.Time
}

// NewUser 创建一个新的用户实例
func NewUser(uuid string, name string) (*User, error) {
	if uuid == "" {
		return nil, errors.New("空的用户uuid")
	}
	if name == "" {
		return nil, errors.New("空的用户名")
	}
	return &User{
		uuid:      uuid,
		name:      name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func UnmarshalUserFromDatabase(
	uuid string,
	name string,
	age uint16,
	gender uint16,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	user, err := NewUser(uuid, name)
	if err != nil {
		return nil, err
	}
	user.age = age
	user.gender = gender
	user.createdAt = createdAt
	user.updatedAt = updatedAt
	return user, nil
}

func (u User) UUID() string {
	return u.uuid
}

func (u User) Name() string {
	return u.name
}

func (u User) Age() uint16 {
	return u.age
}
func (u User) Gender() uint16 {
	return u.gender
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) ChangeUserName(userName string) error {
	if userName == u.name {
		return nil
	}
	if userName == "" {
		return errors.New("can't set user name to empty string")
	}
	u.name = userName
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangeGender(gender uint16) error {
	if !(gender == 0 || gender == 1 || gender == 2) {
		return errors.New("gender must be 0 or 1 or 2")
	}
	u.gender = gender
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangeAge(age uint16) error {
	if age >= 150 {
		return errors.New("age must be less than 150")
	}
	u.gender = age
	u.updatedAt = time.Now()
	return nil
}

func UserIsEqual(u1, u2 *User) bool {
	return u1.uuid == u2.uuid
}
