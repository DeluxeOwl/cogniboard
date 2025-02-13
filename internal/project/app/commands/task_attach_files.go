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
	repo           project.TaskRepository
	fileStorage    project.FileStorage
	embeddings     project.EmbeddingStorage
	imageDescriber project.ImageDescriber
}

func NewAttachFilesToTaskHandler(
	repo project.TaskRepository,
	logger *slog.Logger,
	fileStorage project.FileStorage,
	embeddings project.EmbeddingStorage,
	imageDescriber project.ImageDescriber,
) AttachFilesToTaskHandler {
	return decorator.ApplyCommandDecorators(
		&attachFilesToTaskHandler{
			repo:           repo,
			fileStorage:    fileStorage,
			embeddings:     embeddings,
			imageDescriber: imageDescriber,
		},
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

		err := h.processFile(ctx, cmd.TaskID, &snap, &buf)
		if err != nil {
			return fmt.Errorf("process file: %w", err)
		}

		if err := h.fileStorage.Store(ctx, cmd.TaskID, snap.Name, bytes.NewReader(buf.Bytes())); err != nil {
			return fmt.Errorf("save file: %w", err)
		}
	}

	return h.repo.AddFiles(ctx, cmd.TaskID, files)
}

func (h *attachFilesToTaskHandler) processFile(
	ctx context.Context,
	taskID project.TaskID,
	snap *project.FileSnapshot,
	buf *bytes.Buffer,
) error {
	if !h.shouldCreateEmbeddings(snap.MimeType) {
		fmt.Printf("%s not supported", snap.MimeType) // TODO: this would be better logged somewhere
		return nil
	}

	if isImage(snap.MimeType) {
		file, err := h.fileStorage.Get(ctx, taskID, snap.Name)
		if err != nil {
			return fmt.Errorf("get file: %w", err)
		}

		description, err := h.imageDescriber.DescribeImage(ctx, file)
		if err != nil {
			return fmt.Errorf("describe file: %w", err)
		}
		fmt.Println(description)
		return nil
	}

	if err := h.embeddings.AddDocuments(ctx, []project.Document{
		{
			ID:      snap.ID,
			Name:    snap.Name,
			Content: buf.String(),
			TaskID:  taskID,
		},
	}); err != nil {
		return fmt.Errorf("add documents: %w", err)
	}

	return nil
}

func (h *attachFilesToTaskHandler) shouldCreateEmbeddings(mimeType string) bool {
	switch mimeType {
	case "text/csv", "text/markdown":
		return true
	default:
		if isImage(mimeType) {
			return true
		}
	}
	return false
}

func isImage(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}
