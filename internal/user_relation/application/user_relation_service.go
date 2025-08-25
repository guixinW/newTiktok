package application

import (
	"log/slog"
	"newTiktoken/internal/user_relation/domain/repository"
)

// UserRelationApplicationService provides user relation-related application services.
type UserRelationApplicationService struct {
	userRelationRepo repository.UserRelationRepository
	logger           *slog.Logger
}

// NewUserRelationApplicationService creates a new UserRelationApplicationService.
func NewUserRelationApplicationService(userRelationRepo repository.UserRelationRepository, logger *slog.Logger) *UserRelationApplicationService {
	return &UserRelationApplicationService{userRelationRepo: userRelationRepo, logger: logger}
}
