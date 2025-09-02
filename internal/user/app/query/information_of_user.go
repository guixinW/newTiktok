package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"newTiktoken/internal/common/auth"
	"newTiktoken/internal/common/decorator"
)

type InformationOfUser struct {
	User auth.User
}

type InformationOfUserHandler decorator.QueryHandler[InformationOfUser, *User]

type InformationOfUserReadModel interface {
	FindInformationOfUser(ctx context.Context, userUUID string) (*User, error)
}
type informationOfUserHandler struct {
	readModel InformationOfUserReadModel
}

func (h informationOfUserHandler) Handle(ctx context.Context, query InformationOfUser) (usr *User, err error) {
	return h.readModel.FindInformationOfUser(ctx, query.User.UUID)
}

func NewInformationForUserHandler(
	readModel InformationOfUserReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) InformationOfUserHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return decorator.ApplyQueryDecorators[InformationOfUser, *User](
		informationOfUserHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}
