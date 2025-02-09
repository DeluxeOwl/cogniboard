package adapters

import (
	"context"
	"fmt"
	"io"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

type fileStorage struct {
	bucket *blob.Bucket
}

func NewFileStorage(ctx context.Context, dir string) (project.FileStorage, error) {
	bucket, err := blob.OpenBucket(ctx, fmt.Sprintf("file://%s", dir))
	if err != nil {
		return nil, fmt.Errorf("failed to open bucket: %w", err)
	}

	return &fileStorage{bucket: bucket}, nil
}

func (s *fileStorage) Store(ctx context.Context, taskID project.TaskID, name string, content io.Reader) error {
	key := s.buildKey(taskID, name)
	w, err := s.bucket.NewWriter(ctx, key, nil)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}
	defer w.Close()

	if _, err := io.Copy(w, content); err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	return nil
}

func (s *fileStorage) Get(ctx context.Context, taskID project.TaskID, name string) (io.ReadCloser, error) {
	iter := s.bucket.List(&blob.ListOptions{Prefix: string(taskID)})
	obj, err := iter.Next(ctx)
	if err == io.EOF {
		return nil, fmt.Errorf("file not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return s.bucket.NewReader(ctx, obj.Key, nil)
}

func (s *fileStorage) Delete(ctx context.Context, taskID project.TaskID, name string) error {
	key := s.buildKey(taskID, name)
	if err := s.bucket.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *fileStorage) buildKey(taskID project.TaskID, name string) string {
	return fmt.Sprintf("%s_%s", taskID, name)
}
