package adapters

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	userDomain "newTiktoken/internal/user/domain/user"
	"time"
)

// GormUser 是用于 GORM 的用户模型
type mysqlUser struct {
	ID             uint64
	UserUUID       string
	Username       string
	Age            sql.NullInt16
	Gender         sql.NullInt16
	FollowingCount uint64
	FollowerCount  uint64
	TotalFavorite  uint64
	WorkCount      uint64
	FavoriteCount  uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MySQLUserRepository struct {
	db *sql.DB
}

func (m MySQLUserRepository) marshalDomainUser(user *userDomain.User) mysqlUser {
	mysqlUser := mysqlUser{
		UserUUID:  user.UUID(),
		Username:  user.Name(),
		Age:       sql.NullInt16{Int16: int16(user.Age()), Valid: true},
		Gender:    sql.NullInt16{Int16: int16(user.Gender()), Valid: true},
		UpdatedAt: user.UpdatedAt(),
		CreatedAt: user.CreatedAt(),
	}
	return mysqlUser
}

func (m MySQLUserRepository) unmarshalUser(user *mysqlUser) (*userDomain.User, error) {
	return userDomain.UnmarshalUserFromDatabase(
		user.UserUUID,
		user.Username,
		uint16(user.Gender.Int16),
		uint16(user.Age.Int16),
		user.CreatedAt,
		user.UpdatedAt,
	)
}

// NewMySQLUserRepository 创建一个新的 MySQL 用户仓库实例
func NewMySQLUserRepository(db *sql.DB) (userDomain.Repository, error) {
	// 自动迁移模式，确保表结构是最新的
	return &MySQLUserRepository{
		db: db,
	}, nil
}

// UpdateUser 保存用户 (创建或更新)
func (m MySQLUserRepository) UpdateUser(ctx context.Context, userUUID string, updateFn func(
	ctx context.Context,
	user *userDomain.User,
) (*userDomain.User, error)) (err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	row := tx.QueryRowContext(ctx, "SELECT user_uuid, user_name, age, gender, created_at, updated_at FROM users WHERE user_uuid = ? FOR UPDATE", userUUID)

	var foundUser mysqlUser
	if err = row.Scan(&foundUser.UserUUID, &foundUser.Username, &foundUser.Age, &foundUser.Gender, &foundUser.CreatedAt, &foundUser.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrapf(err, "user with uuid %s not found for update", userUUID)
		}
		return errors.Wrap(err, "failed to scan user for update")
	}

	domainUser, err := m.unmarshalUser(&foundUser)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal user for update")
	}

	updatedDomainUser, err := updateFn(ctx, domainUser)
	if err != nil {
		return errors.Wrap(err, "update function failed")
	}

	updateQuery := "UPDATE users SET user_name = ?, age = ?, gender = ?, updated_at = ? WHERE user_uuid = ?"
	_, err = tx.ExecContext(ctx, updateQuery, updatedDomainUser.Name(), updatedDomainUser.Age(), updatedDomainUser.Gender(), time.Now().UTC(), updatedDomainUser.UUID())
	if err != nil {
		return errors.Wrap(err, "failed to execute user update")
	}

	return nil
}

// AddUser 添加用户
func (m MySQLUserRepository) AddUser(ctx context.Context, user *userDomain.User) error {
	insertQuery := "INSERT INTO users (user_uuid, user_name, age, gender, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	now := time.Now().UTC()
	_, err := m.db.ExecContext(ctx, insertQuery, user.UUID(), user.Name(), user.Age(), user.Gender(), now, now)
	if err != nil {
		return errors.Wrapf(err, "failed to insert user %s", user.UUID())
	}
	return nil
}

// GetUser 根据用户UUID查找用户
func (m MySQLUserRepository) GetUser(ctx context.Context, userUUID string) (*userDomain.User, error) {
	selectQuery := "SELECT user_uuid, user_name, age, gender, created_at, updated_at FROM users WHERE user_uuid = ?"
	row := m.db.QueryRowContext(ctx, selectQuery, userUUID)

	var dbUser mysqlUser
	if err := row.Scan(&dbUser.UserUUID, &dbUser.Username, &dbUser.Age, &dbUser.Gender, &dbUser.CreatedAt, &dbUser.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to find user by uuid %s", userUUID)
	}

	return m.unmarshalUser(&dbUser)
}
