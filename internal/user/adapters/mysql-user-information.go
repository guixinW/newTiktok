package adapters

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"newTiktoken/internal/user/app/query"
)

type MySQLUserFinder struct {
	db *sql.DB
}

func NewMySQLUserFinder(db *sql.DB) (*MySQLUserFinder, error) {
	return &MySQLUserFinder{
		db: db,
	}, nil
}

func (m MySQLUserFinder) FindInformationOfUser(ctx context.Context, userUUID string) (*query.User, error) {
	selectQuery := `
        SELECT
            user_uuid, user_name, age, gender,
            following_count, follower_count, total_favorite, work_count, favorite_count,
            created_at, updated_at
        FROM users
        WHERE user_uuid = ?`
	row := m.db.QueryRowContext(ctx, selectQuery, userUUID)

	var userDTO query.User
	var age sql.NullInt16
	var gender sql.NullInt16

	err := row.Scan(
		&userDTO.UUID,
		&userDTO.Name,
		&age,
		&gender,
		&userDTO.FollowingCount,
		&userDTO.FollowerCount,
		&userDTO.TotalFavorite,
		&userDTO.WorkCount,
		&userDTO.FavoriteCount,
		&userDTO.CreatedAt,
		&userDTO.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to scan user information for %s", userUUID)
	}
	return &userDTO, nil
}
