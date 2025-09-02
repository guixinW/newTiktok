package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"newTiktoken/internal/common/decorator"
	"newTiktoken/internal/common/logs"
	"newTiktoken/internal/user/domain/user"
	"time"
)

type UpdateUser struct {
	UUID      string
	Name      string
	Gender    uint16
	Age       uint16
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUserHandler decorator.CommandHandler[UpdateUser]

type updateUserHandler struct {
	repo user.Repository
}

func (c updateUserHandler) Handle(ctx context.Context, cmd UpdateUser) (err error) {
	defer func() {
		logs.LogCommandExecution("UpdateUser", cmd, err)
	}()
	if err := c.repo.UpdateUser(ctx, cmd.UUID, func(ctx context.Context, user *user.User) (*user.User, error) {
		err := user.ChangeUserName(cmd.Name)
		if err != nil {
			return nil, err
		}
		err = user.ChangeGender(cmd.Gender)
		if err != nil {
			return nil, err
		}
		err = user.ChangeAge(cmd.Age)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func NewUpdateUserHandler(repo user.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient) UpdateUserHandler {
	if repo == nil {
		panic("nil repo")
	}
	return decorator.ApplyCommandDecorators[UpdateUser](
		updateUserHandler{repo: repo},
		logger,
		metricsClient,
	)
}
