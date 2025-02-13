package adapters

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockEmbedding is a simple embedding function that creates a vector based on keyword matching
// This is just for testing purposes and not meant to be used in production
func mockEmbedding(_ context.Context, text string) ([]float32, error) {
	// Create a fixed size vector of 4 dimensions to represent presence of key terms
	vector := make([]float32, 4)

	// Simple keyword matching for our test cases
	keywords := map[string]int{
		"sky":    0, // First dimension
		"blue":   0, // First dimension
		"green":  1, // Second dimension
		"water":  2, // Third dimension
		"boils":  2, // Third dimension
		"degree": 3, // Fourth dimension
	}

	// Convert text to lowercase for case-insensitive matching
	text = strings.ToLower(text)

	// Count keyword occurrences
	for keyword, dimension := range keywords {
		if strings.Contains(text, keyword) {
			vector[dimension] += 1.0
		}
	}

	// Normalize the vector
	var sum float32
	for _, v := range vector {
		sum += v * v
	}
	if sum > 0 {
		magnitude := float32(math.Sqrt(float64(sum)))
		for i := range vector {
			vector[i] /= magnitude
		}
	}

	return vector, nil
}

func TestChromemDB(t *testing.T) {
	// Create a temporary directory for the test database
	tempDir, err := os.MkdirTemp("", "chromem-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")

	// Initialize ChromemDB with mock embedding function
	db, err := NewChromemDB(dbPath, "", "test-collection", mockEmbedding)
	require.NoError(t, err)

	ctx := context.Background()

	// Test documents
	docs := []project.Document{
		{
			ID:      "1",
			TaskID:  "task1",
			Name:    "doc1",
			Content: "The sky is blue because of Rayleigh scattering.",
		},
		{
			ID:      "2",
			TaskID:  "task1",
			Name:    "doc2",
			Content: "Leaves are green because chlorophyll absorbs red and blue light.",
		},
		{
			ID:      "3",
			TaskID:  "task2",
			Name:    "doc3",
			Content: "Water boils at 100 degrees Celsius at sea level.",
		},
	}

	// Test adding documents
	err = db.AddDocuments(ctx, docs)
	require.NoError(t, err)

	// Test searching documents for a specific task
	t.Run("search documents for task", func(t *testing.T) {
		result, err := db.SearchDocumentsForTask(ctx, "task1", "Why is the sky blue?")
		require.NoError(t, err)
		assert.Equal(t, "1", result.ID)
		assert.Equal(t, "doc1", result.Name)
		assert.Contains(t, result.Content, "sky is blue")
	})

	// Test searching all documents
	t.Run("search all documents", func(t *testing.T) {
		results, err := db.SearchAllDocuments(ctx, "What temperature does water boil at?")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "3", results[0].ID)
		assert.Equal(t, "doc3", results[0].Name)
		assert.Contains(t, results[0].Content, "Water boils")
	})

	// Test searching with no results
	t.Run("search with no results", func(t *testing.T) {
		_, err := db.SearchDocumentsForTask(ctx, "nonexistent-task", "query")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no results found")
	})
}
