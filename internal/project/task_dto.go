package project

import (
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent"
	"github.com/danielgtaylor/huma/v2"
)

// In DTOs - for input adapters: e.g REST api

// note: huma doesnt play well with struct embedding
type CreateTask struct {
	Title        string          `form:"title" doc:"Task's name" minLength:"1" maxLength:"50" required:"true"`
	Description  string          `form:"description" doc:"Task's description"`
	DueDate      time.Time       `form:"due_date" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName string          `form:"assignee_name" doc:"Task's asignee (if any)"`
	Files        []huma.FormFile `form:"files"`
}

type EditTask struct {
	Title        string          `form:"title" doc:"Task's name" minLength:"1" maxLength:"50" required:"true"`
	Description  string          `form:"description" doc:"Task's description"`
	DueDate      time.Time       `form:"due_date" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName string          `form:"assignee_name" doc:"Task's asignee (if any)"`
	Files        []huma.FormFile `form:"files"`
}

type ChangeTaskStatus struct {
	Status string `json:"status" doc:"New status for the task" minLength:"1"`
}

type ListTasks struct {
	Tasks []TaskDTO `json:"tasks"`
}

type TaskDTO struct {
	TaskSnapshot
	Files []FileDTO `json:"files"`
}

type FileDTO struct {
	FileSnapshot
}

func ConvertTaskToDTO(task *Task) TaskDTO {
	dto := TaskDTO{
		TaskSnapshot: *task.GetSnapshot(),
		Files:        make([]FileDTO, len(task.files)),
	}

	for i := range task.files {
		dto.Files[i] = ConvertFileToDTO(&task.files[i])
	}

	return dto
}

func ConvertFileToDTO(file *File) FileDTO {
	return FileDTO{
		file.GetSnapshot(),
	}
}

// Out DTOs - for output adapters: e.g. postgres
// DB adapters use GetSnapshot

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
