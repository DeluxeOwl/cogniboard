package adapters

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileStorage(t *testing.T) {
	t.Run("successfully creates file storage", func(t *testing.T) {
		// Create temp directory
		tempDir := t.TempDir()

		// Create new file storage
		ctx := context.Background()
		storage, err := NewFileStorage(ctx, tempDir)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, storage)
	})

	t.Run("fails with invalid directory", func(t *testing.T) {
		ctx := context.Background()
		storage, err := NewFileStorage(ctx, "\x00invalid") // Invalid path character

		require.Error(t, err)
		require.Nil(t, storage)
	})
}

func TestFileStorage_Store(t *testing.T) {
	// Setup
	ctx := context.Background()
	tempDir := t.TempDir()
	storage, err := NewFileStorage(ctx, tempDir)
	require.NoError(t, err)

	t.Run("successfully stores file", func(t *testing.T) {
		// Prepare test data
		taskID := project.TaskID("task123")
		fileName := "test.txt"
		content := "test content"
		reader := strings.NewReader(content)

		// Store file
		err := storage.Store(ctx, taskID, fileName, reader)
		require.NoError(t, err)

		// Verify file exists
		filePath := filepath.Join(tempDir, "task123_test.txt")
		_, err = os.Stat(filePath)
		require.NoError(t, err)

		// Verify content
		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, content, string(data))
	})
}

func TestFileStorage_Get(t *testing.T) {
	// Setup
	ctx := context.Background()
	tempDir := t.TempDir()
	storage, err := NewFileStorage(ctx, tempDir)
	require.NoError(t, err)

	t.Run("successfully gets file", func(t *testing.T) {
		// Store a file first
		taskID := project.TaskID("task123")
		fileName := "test.txt"
		content := "test content"
		err := storage.Store(ctx, taskID, fileName, strings.NewReader(content))
		require.NoError(t, err)

		// Get the file
		reader, err := storage.Get(ctx, taskID, fileName)
		require.NoError(t, err)
		defer reader.Close()

		// Read and verify content
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, content, string(data))
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		taskID := project.TaskID("nonexistent")
		fileName := "nonexistent.txt"

		reader, err := storage.Get(ctx, taskID, fileName)
		require.Error(t, err)
		require.Nil(t, reader)
		assert.Contains(t, err.Error(), "file not found")
	})
}

func TestFileStorage_buildKey(t *testing.T) {
	storage := &fileStorage{}

	testCases := []struct {
		name     string
		taskID   project.TaskID
		fileName string
		want     string
	}{
		{
			name:     "simple file name",
			taskID:   project.TaskID("task123"),
			fileName: "test.txt",
			want:     "task123_test.txt",
		},
		{
			name:     "file name without extension",
			taskID:   project.TaskID("task123"),
			fileName: "test",
			want:     "task123_test",
		},
		{
			name:     "file name with multiple dots",
			taskID:   project.TaskID("task123"),
			fileName: "test.backup.txt",
			want:     "task123_test.backup.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := storage.buildKey(tc.taskID, tc.fileName)
			assert.Equal(t, tc.want, got)
		})
	}
}
