package command

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"newTiktoken/internal/common/decorator"
	"newTiktoken/internal/common/logs"
	"newTiktoken/internal/user/domain/user"
	"time"
)

type CreateUser struct {
	UUID      string
	Name      string
	Age       uint16
	Gender    uint16
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserHandler decorator.CommandHandler[CreateUser]

type createUserHandler struct {
	repo user.Repository
}

func (c createUserHandler) Handle(ctx context.Context, cmd CreateUser) (err error) {
	defer func() {
		logs.LogCommandExecution("CreateUser", cmd, err)
	}()
	existingUser, err := c.repo.GetUser(ctx, cmd.UUID)
	if err != nil {
		return errors.Wrapf(err, "check %s is existed faild", cmd.UUID)
	}
	if existingUser != nil {
		return nil
	}
	usr, err := user.NewUser(cmd.UUID, cmd.Name)
	if err != nil {
		return err
	}
	if err := c.repo.AddUser(ctx, usr); err != nil {
		return err
	}
	return nil
}

func NewCreateUserHandler(repo user.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient) CreateUserHandler {
	if repo == nil {
		panic("nil repo")
	}
	return decorator.ApplyCommandDecorators[CreateUser](
		createUserHandler{repo: repo},
		logger,
		metricsClient,
	)
}
