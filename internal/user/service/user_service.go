package service

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"newTiktoken/internal/common/metrics"
	"newTiktoken/internal/user/adapters"
	"newTiktoken/internal/user/app"
	"newTiktoken/internal/user/app/command"
	"newTiktoken/internal/user/app/query"
	"os"
)

func NewApplication(ctx context.Context) app.Application {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	userRepository, err := adapters.NewMySQLUserRepository(db)
	if err != nil {
		panic(err)
	}
	userFinder, err := adapters.NewMySQLUserFinder(db)
	if err != nil {
		panic(err)
	}
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	if err != nil {
		panic(err)
	}
	return app.Application{
		Commands: app.Commands{
			UpdateUser: command.NewUpdateUserHandler(userRepository, logger, metricsClient),
			CreateUser: command.NewCreateUserHandler(userRepository, logger, metricsClient),
		},
		Queries: app.Queries{
			InformationOfUser: query.NewInformationForUserHandler(userFinder, logger, metricsClient),
		},
	}
}
