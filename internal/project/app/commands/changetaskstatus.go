package commands

import (
	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type ChangeTaskStatus struct {
	TaskID string
	Status string
}

type ChangeTaskStatusHandler decorator.CommandHandler[ChangeTaskStatus]

type changeTaskStatusHandler struct {
	repo project.TaskRepository
}
