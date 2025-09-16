package adapters

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	userRelationDomain "newTiktoken/internal/user-relation/domain"
	"time"
)

type mysqlUserRelation struct {
	ID               uint64
	ActivePartyUUID  string
	PassivePartyUUID string
	Status           int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type MySQLUserRelationRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) (userRelationDomain.Repository, error) {
	return &MySQLUserRelationRepository{
		db: db,
	}, nil
}

func (m MySQLUserRelationRepository) UpdateRelation(
	ctx context.Context,
	ActivePartyUUID string,
	PassivePartyUUID string,
	updateFn func(ctx context.Context, userRelation *userRelationDomain.UserRelation) (*userRelationDomain.UserRelation, error)) error {
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
	const findAndLockQuery = `
        SELECT active_party_uuid, passive_party_uuid, status, created_at, updated_at
        FROM user_relations
        WHERE active_party_uuid = ? and passive_party_uuid = ?
        FOR UPDATE`

	row := tx.QueryRowContext(ctx, findAndLockQuery, ActivePartyUUID, PassivePartyUUID)
	var foundRelation mysqlUserRelation
	err = row.Scan(
		&foundRelation.ActivePartyUUID,
		&foundRelation.PassivePartyUUID,
		&foundRelation.Status,
		&foundRelation.CreatedAt,
		&foundRelation.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrapf(err, "user relation between uuid %s and uuid %s not found for update", ActivePartyUUID, PassivePartyUUID)
		}
		return errors.Wrap(err, "failed to scan user relation for update")
	}
	domainUserRelation, err := m.unmarshalUser(&foundRelation)
	updatedRelation, err := updateFn(ctx, domainUserRelation)
	if err != nil {
		return errors.Wrap(err, "failed to update user relation")
	}
	updateQuery := "UPDATE user_relations SET status = ? WHERE active_party_uuid = ? AND passive_party_uuid = ?"
	_, err = tx.ExecContext(ctx, updateQuery,
		updatedRelation.Status,
		updatedRelation.ActivePartyUUID,
		updatedRelation.PassivePartyUUID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update user relation")
	}
	return nil
}

func (m MySQLUserRelationRepository) GetRelation(ctx context.Context, ActivePartyUUID, PassivePartyUUID string) (*userRelationDomain.UserRelation, error) {
	const query = `
        SELECT id, active_party_uuid, passive_party_uuid, status, created_at, updated_at
        FROM user_relations
        WHERE active_party_uuid = ? AND passive_party_uuid = ?`

	row := m.db.QueryRowContext(ctx, query, ActivePartyUUID, PassivePartyUUID)

	var relation mysqlUserRelation
	err := row.Scan(
		&relation.ID,
		&relation.ActivePartyUUID,
		&relation.PassivePartyUUID,
		&relation.Status,
		&relation.CreatedAt,
		&relation.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "relation not found")
		}
		return nil, errors.Wrap(err, "failed to scan user relation")
	}

	return m.unmarshalUser(&relation)
}

func (m MySQLUserRelationRepository) AddRelation(ctx context.Context, ActivePartyUUID, PassivePartyUUID string) error {
	const defaultFollowStatus = 0
	const query = `
        INSERT INTO user_relations (active_party_uuid, passive_party_uuid, status)
        VALUES (?, ?, ?)`
	_, err := m.db.ExecContext(ctx, query, ActivePartyUUID, PassivePartyUUID, defaultFollowStatus)
	if err != nil {
		return errors.Wrap(err, "failed to add relation")
	}

	return nil
}

func (m MySQLUserRelationRepository) unmarshalUser(relation *mysqlUserRelation) (*userRelationDomain.UserRelation, error) {
	relationActionType, err := userRelationDomain.NewRelationTypeFromInt(relation.Status)
	if err != nil {
		return nil, err
	}
	return userRelationDomain.UnmarshalUserRelationFromDatabase(
		relation.ActivePartyUUID,
		relation.PassivePartyUUID,
		relationActionType,
		relation.CreatedAt,
		relation.UpdatedAt,
	)
}
