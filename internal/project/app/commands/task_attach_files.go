package commands

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type FileToUpload struct {
	Metadata project.File
	Content  io.Reader
}

type AttachFilesToTask struct {
	TaskID project.TaskID
	Files  []FileToUpload
}

type AttachFilesToTaskHandler decorator.CommandHandler[AttachFilesToTask]

type attachFilesToTaskHandler struct {
	repo        project.TaskRepository
	fileStorage project.FileStorage
	embeddings  project.EmbeddingStorage
}

func NewAttachFilesToTaskHandler(repo project.TaskRepository, logger *slog.Logger, fileStorage project.FileStorage, embeddings project.EmbeddingStorage) AttachFilesToTaskHandler {
	return decorator.ApplyCommandDecorators(
		&attachFilesToTaskHandler{repo: repo, fileStorage: fileStorage, embeddings: embeddings},
		logger,
	)
}

func (h *attachFilesToTaskHandler) Handle(ctx context.Context, cmd AttachFilesToTask) error {
	if len(cmd.Files) == 0 {
		return nil
	}

	files := make([]project.File, len(cmd.Files))
	for i, file := range cmd.Files {
		files[i] = file.Metadata
		err := h.fileStorage.Store(ctx, cmd.TaskID, file.Metadata.GetSnapshot().Name, file.Content)
		if err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}
	}

	return h.repo.AddFiles(ctx, cmd.TaskID, files)
}
