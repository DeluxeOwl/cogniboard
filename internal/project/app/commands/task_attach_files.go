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
	logger         *slog.Logger
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
			logger:         logger,
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

		if err := h.fileStorage.Store(ctx, cmd.TaskID, snap.Name, bytes.NewReader(buf.Bytes())); err != nil {
			return fmt.Errorf("save file: %w", err)
		}

		go func() {
			ctx := context.Background()
			err := h.processFile(ctx, cmd.TaskID, &snap, &buf)
			if err != nil {
				h.logger.Error("file not processed", "err", err)
			}
		}()

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
		h.logger.Warn("file type not supported for processing", "mime_type", snap.MimeType)
		return nil
	}

	content, err := h.getFileContent(ctx, taskID, snap, buf)
	if err != nil {
		return err
	}

	return h.addDocumentEmbedding(ctx, project.Document{
		ID:      snap.ID,
		Name:    snap.Name,
		Content: content,
		TaskID:  taskID,
	})
}

func (h *attachFilesToTaskHandler) getFileContent(
	ctx context.Context,
	taskID project.TaskID,
	snap *project.FileSnapshot,
	buf *bytes.Buffer,
) (string, error) {
	if !isImage(snap.MimeType) {
		return buf.String(), nil
	}

	file, err := h.fileStorage.Get(ctx, taskID, snap.Name)
	if err != nil {
		return "", fmt.Errorf("get file: %w", err)
	}
	defer file.Close()

	description, err := h.imageDescriber.DescribeImage(ctx, file)
	if err != nil {
		return "", fmt.Errorf("describe file: %w", err)
	}
	slog.Info("describe image", "description", description)

	return description, nil
}

func (h *attachFilesToTaskHandler) addDocumentEmbedding(
	ctx context.Context,
	doc project.Document,
) error {
	if err := h.embeddings.AddDocuments(ctx, []project.Document{doc}); err != nil {
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
