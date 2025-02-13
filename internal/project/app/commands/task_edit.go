package commands

import (
	"context"
	"log/slog"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type EditTask struct {
	TaskID       string
	Title        *string
	Description  *string
	DueDate      *time.Time
	AssigneeName *string
	Status       *string
}

type EditTaskHandler decorator.CommandHandler[EditTask]

type editTaskHandler struct {
	repo project.TaskRepository
}

func NewEditTaskHandler(repo project.TaskRepository, logger *slog.Logger) EditTaskHandler {
	return decorator.ApplyCommandDecorators(
		&editTaskHandler{repo: repo},
		logger,
	)
}

func (h *editTaskHandler) Handle(ctx context.Context, cmd EditTask) error {
	return h.repo.UpdateTask(
		ctx,
		project.TaskID(cmd.TaskID),
		func(t *project.Task) (*project.Task, error) {
			var status *project.TaskStatus
			if cmd.Status != nil {
				s := project.TaskStatus(*cmd.Status)
				status = &s
			}

			if err := t.Edit(cmd.Title, cmd.Description, cmd.DueDate, cmd.AssigneeName, status); err != nil {
				return nil, err
			}
			return t, nil
		},
	)
}
