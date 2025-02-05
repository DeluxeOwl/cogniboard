package queries

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type AllTasks struct{}

type AllTasksHandler decorator.QueryHandler[AllTasks, []project.Task]

type allTasksHandler struct {
	tasks AllTasksReadModel
}

type AllTasksReadModel interface {
	AllTasks(ctx context.Context) ([]project.Task, error)
}

func NewAllTasksHandler(repo AllTasksReadModel, logger *slog.Logger) (AllTasksHandler, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}

	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return decorator.ApplyQueryDecorators(
		&allTasksHandler{tasks: repo},
		logger,
	), nil
}
func (h *allTasksHandler) Handle(ctx context.Context, query AllTasks) ([]project.Task, error) {
	return h.tasks.AllTasks(ctx)
}
