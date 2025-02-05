package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type UnassignTask struct {
	TaskID string
}

type UnassignTaskHandler decorator.CommandHandler[UnassignTask]

type unassignTaskHandler struct {
	repo project.TaskRepository
}

func NewUnassignTaskHandler(repo project.TaskRepository, logger *slog.Logger) (UnassignTaskHandler, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return decorator.ApplyCommandDecorators(
		&unassignTaskHandler{repo: repo},
		logger,
	), nil
}

func (h *unassignTaskHandler) Handle(ctx context.Context, cmd UnassignTask) error {
	return h.repo.UpdateTask(ctx, project.TaskID(cmd.TaskID), func(t *project.Task) (*project.Task, error) {
		t.Unassign()
		return t, nil
	})
}
