package commands

import (
	"context"
	"log/slog"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type AttachFilesToTask struct {
	TaskID project.TaskID
	Files  []project.File
}

type AttachFilesToTaskHandler decorator.CommandHandler[AttachFilesToTask]

type attachFilesToTaskHandler struct {
	repo project.TaskRepository
}

func NewAttachFilesToTaskHandler(repo project.TaskRepository, logger *slog.Logger) AttachFilesToTaskHandler {
	return decorator.ApplyCommandDecorators(
		&attachFilesToTaskHandler{repo: repo},
		logger,
	)
}

func (h *attachFilesToTaskHandler) Handle(ctx context.Context, cmd AttachFilesToTask) error {
	if len(cmd.Files) == 0 {
		return nil
	}

	return h.repo.AddFiles(ctx, cmd.TaskID, cmd.Files)
}
