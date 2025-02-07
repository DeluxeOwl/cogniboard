package app

import (
	"errors"
	"log/slog"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/commands"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/queries"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateTask       commands.CreateTaskHandler
	ChangeTaskStatus commands.ChangeTaskStatusHandler
	EditTask         commands.EditTaskHandler
}

type Queries struct {
	AllTasks queries.AllTasksHandler
}

// Validation for the arguments injected into command & queries happens here
func New(repo project.TaskRepository, logger *slog.Logger) (*Application, error) {
	if repo == nil {
		return nil, errors.New("repo cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &Application{
		Commands: Commands{
			CreateTask:       commands.NewCreateTaskHandler(repo, logger),
			ChangeTaskStatus: commands.NewChangeStatusHandler(repo, logger),
			EditTask:         commands.NewEditTaskHandler(repo, logger),
		},
		Queries: Queries{
			AllTasks: queries.NewAllTasksHandler(repo, logger),
		},
	}, nil
}
