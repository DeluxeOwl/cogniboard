package adapters

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/philippgille/chromem-go"
)

type ChromemDB struct {
	collection   *chromem.Collection
	taskField    string
	docNameField string
}

var _ project.EmbeddingStorage = &ChromemDB{}

func NewChromemDB(dbPath string, openAIKey string, collectionName string, embeddingFunc chromem.EmbeddingFunc) (*ChromemDB, error) {
	db, err := chromem.NewPersistentDB(dbPath, false)
	if err != nil {
		return nil, err
	}

	c, err := db.GetOrCreateCollection(collectionName, nil, embeddingFunc)
	if err != nil {
		return nil, err
	}
	return &ChromemDB{collection: c, taskField: "task_id", docNameField: "document_name"}, nil
}

func (c *ChromemDB) AddDocuments(ctx context.Context, docs []project.Document) error {
	chromemDocs := make([]chromem.Document, len(docs))

	for i, doc := range docs {
		chromemDocs[i] = chromem.Document{
			ID:      doc.ID,
			Content: doc.Content,
			Metadata: map[string]string{
				c.taskField:    string(doc.TaskID),
				c.docNameField: doc.Name,
			},
		}
	}

	return c.collection.AddDocuments(ctx, chromemDocs, runtime.NumCPU())
}

func (c *ChromemDB) SearchDocumentsForTask(ctx context.Context, taskID project.TaskID, query string) (*project.DocumentSimilarity, error) {
	res, err := c.collection.Query(ctx, query, 1, map[string]string{
		c.taskField: string(taskID),
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("searching documents for task id %q and query %q failed: %w", taskID, query, err)
	}

	if len(res) == 0 {
		return nil, errors.New("no results found")
	}

	return &project.DocumentSimilarity{
		ID:         res[0].ID,
		Name:       res[0].Metadata[c.docNameField],
		Content:    res[0].Content,
		Similarity: res[0].Similarity,
		Metadata:   res[0].Metadata,
	}, nil
}

func (c *ChromemDB) SearchAllDocuments(ctx context.Context, query string) ([]project.DocumentSimilarity, error) {
	res, err := c.collection.Query(ctx, query, 1, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("searching all documents for query %q failed: %w", query, err)
	}

	if len(res) == 0 {
		return nil, errors.New("no results found")
	}

	docs := make([]project.DocumentSimilarity, len(res))
	for i, doc := range res {
		docs[i] = project.DocumentSimilarity{
			ID:         doc.ID,
			Name:       doc.Metadata[c.docNameField],
			Content:    doc.Content,
			Similarity: doc.Similarity,
			Metadata:   doc.Metadata,
		}
	}
	return docs, nil
}
