package project

import (
	"time"
)

// In DTOs - for input adapters: e.g REST api
type InCreateTaskDTO struct {
	Title        string     `json:"title" doc:"Task's name" minLength:"1" maxLength:"50"`
	Description  *string    `json:"description,omitempty" doc:"Task's description"`
	DueDate      *time.Time `json:"due_date,omitempty" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName *string    `json:"assignee_name,omitempty" doc:"Task's asignee (if any)"`
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

type InTasksDTO struct {
	Tasks []InTaskDTO `json:"tasks"`
}

type InEditTaskDTO struct {
	Title        *string    `json:"title,omitempty" doc:"Task's title" maxLength:"50"`
	Description  *string    `json:"description,omitempty" doc:"Task's description"`
	DueDate      *time.Time `json:"due_date,omitempty" doc:"Task's due date" format:"date-time"`
	AssigneeName *string    `json:"assignee_name,omitempty" doc:"Name of the person to assign the task to"`
	Status       *string    `json:"status,omitempty" doc:"Task's status"`
}

type InChangeTaskStatusDTO struct {
	Status string `json:"status" doc:"New status for the task" minLength:"1"`
}

// Out DTOs - for output adapters: e.g. postgres

type OutTaskDTO struct {
	ID           string     `db:"id"`
	Title        string     `db:"title"`
	Description  *string    `db:"description"`
	DueDate      *time.Time `db:"due_date"`
	AssigneeName *string    `db:"assignee"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	CompletedAt  *time.Time `db:"completed_at"`
	Status       string     `db:"status"`
}

func ToOutTaskDTO(t *Task) *OutTaskDTO {
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

	return task, nil
}
