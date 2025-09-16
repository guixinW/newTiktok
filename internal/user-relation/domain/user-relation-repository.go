package domain

import "context"

type Repository interface {
	GetRelation(ctx context.Context, ActivePartyUUID, PassivePartyUUID string) (*UserRelation, error)
	AddRelation(ctx context.Context, ActivePartyUUID, PassivePartyUUID string) error
	UpdateRelation(ctx context.Context, ActivePartyUUID, PassivePartyUUID string, updateFn func(
		ctx context.Context,
		userRelation *UserRelation,
	) (*UserRelation, error)) error
}
