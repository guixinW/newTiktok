package domain

import (
	"fmt"
	commonError "newTiktoken/internal/common/errors"
)

var (
	Follow   = RelationActionType{1}
	Unfollow = RelationActionType{2}
	Block    = RelationActionType{3}
)

type RelationActionType struct {
	action int
}

func NewRelationTypeFromInt(relationActionType int) (RelationActionType, error) {
	switch relationActionType {
	case 1:
		return Follow, nil
	case 2:
		return Unfollow, nil
	case 3:
		return Block, nil
	}

	return RelationActionType{}, commonError.NewIncorrectInputError(
		fmt.Sprintf("invalid '%s' role", relationActionType),
		"invalid-role",
	)
}
