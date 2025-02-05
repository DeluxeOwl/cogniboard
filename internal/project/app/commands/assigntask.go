package commands

import (
	"context"
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

func NewAssignTaskHandler(repo project.TaskRepository, logger *slog.Logger) AssignTaskHandler {
	return decorator.ApplyCommandDecorators(
		&assignTaskHandler{repo: repo},
		logger,
	)
}

func (h *assignTaskHandler) Handle(ctx context.Context, cmd AssignTask) error {
	return h.repo.UpdateTask(ctx, project.TaskID(cmd.TaskID), func(t *project.Task) (*project.Task, error) {
		if err := t.AssignTo(&cmd.AssigneeName); err != nil {
			return nil, err
		}
		return t, nil
	})
}
