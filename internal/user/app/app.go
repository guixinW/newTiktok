package app

import (
	"newTiktoken/internal/user/app/command"
	"newTiktoken/internal/user/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateUser command.CreateUserHandler
	UpdateUser command.UpdateUserHandler
}

type Queries struct {
	InformationOfUser query.InformationOfUserHandler
}
