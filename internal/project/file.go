package project

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"time"
)

type FileStorage interface {
	Store(ctx context.Context, taskID TaskID, name string, content io.Reader) error
	Get(ctx context.Context, taskID TaskID, name string) (io.ReadCloser, error)
}

// File is a value object that represents a file attached to a task
type File struct {
	Name       string
	Size       int64
	MimeType   string
	UploadedAt time.Time
}

var (
	ErrInvalidFileName = fmt.Errorf("invalid file name")
	ErrInvalidFileSize = fmt.Errorf("invalid file size")
)

func NewFile(name string, size int64) (File, error) {
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
		Name:       name,
		Size:       size,
		MimeType:   mimeType,
		UploadedAt: time.Now(),
	}, nil
}
