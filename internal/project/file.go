package project

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type FileStorage interface {
	Store(ctx context.Context, taskID TaskID, name string, content io.Reader) error
	Get(ctx context.Context, taskID TaskID, name string) (io.ReadCloser, error)
}

// File is a value object that represents a file attached to a task
type File struct {
	id         string
	name       string
	size       int64
	mimeType   string
	uploadedAt time.Time
}

var (
	ErrInvalidFileName = fmt.Errorf("invalid file name")
	ErrInvalidFileSize = fmt.Errorf("invalid file size")
)

func NewFile(name string, size int64) (File, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return File{}, err
	}

	if name == "" {
		return File{}, ErrInvalidFileName
	}

	if size <= 0 {
		return File{}, ErrInvalidFileSize
	}

	mimeType := mime.TypeByExtension(filepath.Ext(name))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return File{
		id:         id.String(),
		name:       name,
		size:       size,
		mimeType:   mimeType,
		uploadedAt: time.Now(),
	}, nil
}

func (f *File) GetSnapshot() FileSnapshot {
	return FileSnapshot{
		ID:         f.id,
		Name:       f.name,
		Size:       f.size,
		MimeType:   f.mimeType,
		UploadedAt: f.uploadedAt,
	}
}

type FileSnapshot struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}
