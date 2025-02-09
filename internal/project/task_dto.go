package project

import (
	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent"
)

// Out DTOs - for output adapters: e.g. postgres
// Adapters use GetSnapshot

func UnmarshalTaskFromDB(t *ent.Task) (*Task, error) {
	task, err := NewTask(TaskID(t.ID), t.Title, t.Description, t.DueDate, t.AssigneeName)
	if err != nil {
		return nil, err
	}

	task.createdAt = t.CreatedAt
	task.completedAt = t.CompletedAt
	task.updatedAt = t.UpdatedAt
	task.status = TaskStatus(t.Status)

	files := make([]File, len(t.Edges.Files))
	for i, f := range t.Edges.Files {
		files[i] = UnmarshalFileFromDB(f)
	}
	task.files = files

	return task, nil
}

func UnmarshalFileFromDB(f *ent.File) File {
	return File{
		id:         f.ID,
		name:       f.Name,
		size:       f.Size,
		mimeType:   f.MimeType,
		uploadedAt: f.UploadedAt,
	}
}
