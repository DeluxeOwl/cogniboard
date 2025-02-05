package commands

import (
	"context"
	"log/slog"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type CreateTask struct {
	Title        string
	Description  *string
	DueDate      *time.Time
	AssigneeName *string
}

type CreateTaskHandler decorator.CommandHandler[CreateTask]

type createTaskHandler struct {
	repo project.TaskRepository
}

func NewCreateTaskHandler(repo project.TaskRepository, logger *slog.Logger) CreateTaskHandler {
	return decorator.ApplyCommandDecorators(
		&createTaskHandler{repo: repo},
		logger,
	)
}

func (h *createTaskHandler) Handle(ctx context.Context, cmd CreateTask) error {
	taskID, err := project.NewTaskID()
	if err != nil {
		return err
	}

	task, err := project.NewTask(taskID, cmd.Title, cmd.Description, cmd.DueDate, cmd.AssigneeName)
	if err != nil {
		return err
	}

	if err := h.repo.Create(ctx, task); err != nil {
		return err
	}

	return nil
}
