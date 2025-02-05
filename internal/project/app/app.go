package app

import (
	"github.com/DeluxeOwl/cogniboard/internal/project/app/commands"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/queries"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ChangeTaskStatus commands.ChangeTaskStatusHandler
	AssignTask       commands.AssignTaskHandler
	UnassignTask     commands.UnassignTaskHandler
}

type Queries struct {
	AllTasks queries.AllTasksHandler
}
