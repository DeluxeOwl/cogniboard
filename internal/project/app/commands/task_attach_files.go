package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

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

func NewAttachFilesToTaskHandler(
	repo project.TaskRepository,
	logger *slog.Logger,
	fileStorage project.FileStorage,
	embeddings project.EmbeddingStorage,
) AttachFilesToTaskHandler {
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
		snap := file.Metadata.GetSnapshot()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file.Content); err != nil {
			return fmt.Errorf("copy file content: %w", err)
		}

		if h.shouldCreateEmbeddings(snap.MimeType) {
			err := h.embeddings.AddDocuments(ctx, []project.Document{
				{
					ID:      snap.ID,
					Name:    snap.Name,
					Content: buf.String(),
					TaskID:  cmd.TaskID,
				},
			})
			if err != nil {
				fmt.Println("TODO: failed to create embedding", err)
			}
		}

		if err := h.fileStorage.Store(ctx, cmd.TaskID, snap.Name, bytes.NewReader(buf.Bytes())); err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}
	}

	return h.repo.AddFiles(ctx, cmd.TaskID, files)
}

func (h *attachFilesToTaskHandler) shouldCreateEmbeddings(mimeType string) bool {
	switch mimeType {
	case "text/csv", "text/markdown":
		return true
	default:
		if strings.HasPrefix(mimeType, "image/") {
			fmt.Println("TODO: image types are not supported")
			return false
		}
	}
	return false
}
