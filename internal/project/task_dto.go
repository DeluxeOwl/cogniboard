package project

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// In DTOs - for input adapters: e.g REST api
type InCreateTaskDTO struct {
	Title        string    `form:"title" doc:"Task's name" minLength:"1" maxLength:"50" required:"true"`
	Description  string    `form:"description" doc:"Task's description"`
	DueDate      time.Time `form:"due_date" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName string    `form:"assignee_name" doc:"Task's asignee (if any)"`

	Files []huma.FormFile `form:"files"`
}

type InTasksDTO struct {
	Tasks []InTaskDTO `json:"tasks"`
}

type InTaskDTO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Assignee    *string    `json:"assignee"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Status      string     `json:"status"`
}

func ToInTaskDTO(task *Task) InTaskDTO {
	return InTaskDTO{
		ID:          string(task.id),
		Title:       task.title,
		Description: task.description,
		DueDate:     task.dueDate,
		Assignee:    task.asigneeName,
		CreatedAt:   task.createdAt,
		CompletedAt: task.completedAt,
		UpdatedAt:   task.updatedAt,
		Status:      string(task.status),
	}
}

type InEditTaskDTO struct {
	Title        string    `form:"title" doc:"Task's name" minLength:"1" maxLength:"50" required:"true"`
	Description  string    `form:"description" doc:"Task's description"`
	DueDate      time.Time `form:"due_date" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName string    `form:"assignee_name" doc:"Task's asignee (if any)"`

	Files []huma.FormFile `form:"files"`
}

type InChangeTaskStatusDTO struct {
	Status string `json:"status" doc:"New status for the task" minLength:"1"`
}

// Out DTOs - for output adapters: e.g. postgres

func FileToOutFileDTO(file *File) OutFileDTO {
	return OutFileDTO{
		Name:       file.Name,
		Size:       file.Size,
		MimeType:   file.MimeType,
		UploadedAt: file.UploadedAt,
	}
}

func FileFromOutFileDTO(dto *OutFileDTO) *File {
	return &File{
		Name:       dto.Name,
		Size:       dto.Size,
		MimeType:   dto.MimeType,
		UploadedAt: dto.UploadedAt,
	}
}

type OutFileDTO struct {
	Name       string    `db:"name"`
	Size       int64     `db:"size"`
	MimeType   string    `db:"mime_type"`
	UploadedAt time.Time `db:"uploaded_at"`
}

type OutTaskDTO struct {
	ID           string       `db:"id"`
	Title        string       `db:"title"`
	Description  *string      `db:"description"`
	DueDate      *time.Time   `db:"due_date"`
	AssigneeName *string      `db:"assignee"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	CompletedAt  *time.Time   `db:"completed_at"`
	Status       string       `db:"status"`
	Files        []OutFileDTO `db:"-"` // Using db:"-" as files will be loaded separately
}

func ToOutTaskDTO(t *Task) *OutTaskDTO {
	taskFiles := t.Files()
	outFiles := make([]OutFileDTO, len(taskFiles))

	for i, file := range taskFiles {
		outFiles[i] = FileToOutFileDTO(&file)
	}

	return &OutTaskDTO{
		ID:           string(t.id),
		Title:        t.title,
		Description:  t.description,
		DueDate:      t.dueDate,
		AssigneeName: t.asigneeName,
		CreatedAt:    t.createdAt,
		UpdatedAt:    t.updatedAt,
		CompletedAt:  t.completedAt,
		Status:       string(t.status),
		Files:        outFiles,
	}
}

func FromOutTaskDTO(t *OutTaskDTO) (*Task, error) {
	task, err := NewTask(TaskID(t.ID), t.Title, t.Description, t.DueDate, t.AssigneeName)
	if err != nil {
		return nil, err
	}

	task.createdAt = t.CreatedAt
	task.updatedAt = t.UpdatedAt
	task.completedAt = t.CompletedAt
	task.status = TaskStatus(t.Status)

	for _, fileDTO := range t.Files {
		file := FileFromOutFileDTO(&fileDTO)
		task.files = append(task.files, *file)
	}

	return task, nil
}
