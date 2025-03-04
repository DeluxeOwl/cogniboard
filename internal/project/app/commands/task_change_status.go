package commands

import (
	"context"
	"log/slog"

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

func NewChangeStatusHandler(
	repo project.TaskRepository,
	logger *slog.Logger,
) ChangeTaskStatusHandler {
	return decorator.ApplyCommandDecorators(
		&changeTaskStatusHandler{repo: repo},
		logger,
	)
}

func (h *changeTaskStatusHandler) Handle(ctx context.Context, cmd ChangeTaskStatus) error {
	if err := h.repo.UpdateTask(ctx, project.TaskID(cmd.TaskID), func(t *project.Task) (*project.Task, error) {
		err := t.ChangeStatus(project.TaskStatus(cmd.Status))
		if err != nil {
			return nil, err
		}
		return t, nil
	}); err != nil {
		return err
	}

	return nil
}
