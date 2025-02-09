package project

import (
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent"
	"github.com/DeluxeOwl/cogniboard/internal/postgres/ent/task"
	"github.com/danielgtaylor/huma/v2"
)

// In DTOs - for input adapters: e.g REST api
type InCreateTaskDTO struct {
	Title        string          `form:"title" doc:"Task's name" minLength:"1" maxLength:"50" required:"true"`
	Description  string          `form:"description" doc:"Task's description"`
	DueDate      time.Time       `form:"due_date" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName string          `form:"assignee_name" doc:"Task's asignee (if any)"`
	Files        []huma.FormFile `form:"files"`
}

type InTasksDTO struct {
	Tasks []InTaskDTO `json:"tasks"`
}

type InTaskDTO struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description *string     `json:"description"`
	DueDate     *time.Time  `json:"due_date"`
	Assignee    *string     `json:"assignee"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CompletedAt *time.Time  `json:"completed_at"`
	Status      string      `json:"status"`
	Files       []InFileDTO `json:"files"`
}

type InFileDTO struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func ToInTaskDTO(task *Task) InTaskDTO {
	return InTaskDTO{
		ID:          string(task.id),
		Title:       task.title,
		Description: task.description,
		DueDate:     task.dueDate,
		Assignee:    task.asigneeName,
		CreatedAt:   task.createdAt,
		UpdatedAt:   task.updatedAt,
		CompletedAt: task.completedAt,
		Status:      string(task.status),
		Files:       ToInFileDTOArray(task.files),
	}
}

func ToInFileDTOArray(files []File) []InFileDTO {
	inFileDTOs := make([]InFileDTO, len(files))

	for i, file := range files {
		inFileDTOs[i] = InFileDTO{
			Name:       file.Name,
			Size:       file.Size,
			MimeType:   file.MimeType,
			UploadedAt: file.UploadedAt,
		}
	}

	return inFileDTOs
}

type InEditTaskDTO struct {
	Title        string          `form:"title" doc:"Task's name" minLength:"1" maxLength:"50" required:"true"`
	Description  string          `form:"description" doc:"Task's description"`
	DueDate      time.Time       `form:"due_date" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName string          `form:"assignee_name" doc:"Task's asignee (if any)"`
	Files        []huma.FormFile `form:"files"`
}

type InChangeTaskStatusDTO struct {
	Status string `json:"status" doc:"New status for the task" minLength:"1"`
}

// Out DTOs - for output adapters: e.g. postgres

func EntToTask(t *ent.Task) (*Task, error) {
	task, err := NewTask(TaskID(t.ID), t.Title, t.Description, t.DueDate, t.AssigneeName)
	if err != nil {
		return nil, err
	}

	task.createdAt = t.CreatedAt
	task.completedAt = t.CompletedAt
	task.updatedAt = t.UpdatedAt
	task.status = TaskStatus(t.Status)

	task.files = EntFilesToTaskFiles(t.Edges.Files)

	return task, nil
}

func TaskToEnt(t *Task) *ent.Task {
	return &ent.Task{
		ID:           string(t.id),
		Title:        t.title,
		Description:  t.description,
		AssigneeName: t.asigneeName,
		DueDate:      t.dueDate,
		CreatedAt:    t.createdAt,
		UpdatedAt:    t.updatedAt,
		CompletedAt:  t.completedAt,
		Status:       task.Status(t.status),
		Edges: ent.TaskEdges{
			Files: TaskFilesToEntFiles(t.files),
		},
	}
}

func TaskFilesToEntFiles(taskFiles []File) []*ent.File {
	files := make([]*ent.File, len(taskFiles))
	for i, f := range taskFiles {
		files[i] = &ent.File{
			ID:         f.ID,
			Name:       f.Name,
			Size:       f.Size,
			MimeType:   f.MimeType,
			UploadedAt: f.UploadedAt,
		}
	}
	return files
}

func EntFilesToTaskFiles(entFiles []*ent.File) []File {
	files := make([]File, len(entFiles))
	for i, f := range entFiles {
		files[i] = File{
			ID:         f.ID,
			Name:       f.Name,
			Size:       f.Size,
			MimeType:   f.MimeType,
			UploadedAt: f.UploadedAt,
		}
	}
	return files
}
