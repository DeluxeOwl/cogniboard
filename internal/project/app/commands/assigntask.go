package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type AssignTask struct {
	TaskID       string
	AssigneeName string
}

type AssignTaskHandler decorator.CommandHandler[AssignTask]

type assignTaskHandler struct {
	repo project.TaskRepository
}

func NewAssignTaskHandler(repo project.TaskRepository, logger *slog.Logger) (AssignTaskHandler, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return decorator.ApplyCommandDecorators(
		&assignTaskHandler{repo: repo},
		logger,
	), nil
}

func (h *assignTaskHandler) Handle(ctx context.Context, cmd AssignTask) error {
	return h.repo.UpdateTask(ctx, project.TaskID(cmd.TaskID), func(t *project.Task) (*project.Task, error) {
		if err := t.AssignTo(&cmd.AssigneeName); err != nil {
			return nil, err
		}
		return t, nil
	})
}
