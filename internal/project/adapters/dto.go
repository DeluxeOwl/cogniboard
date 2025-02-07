package adapters

import "time"

// CreateTaskDTO represents the input for creating a task
type CreateTaskDTO struct {
	Title        string     `json:"title" doc:"Task's name" minLength:"1" maxLength:"50"`
	Description  *string    `json:"description,omitempty" doc:"Task's description"`
	DueDate      *time.Time `json:"due_date,omitempty" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName *string    `json:"assignee_name,omitempty" doc:"Task's asignee (if any)"`
}

// TaskDTO represents a task in the response
type TaskDTO struct {
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

// TasksDTO represents a collection of tasks
type TasksDTO struct {
	Tasks []TaskDTO `json:"tasks"`
}

// EditTaskDTO represents the input for editing a task
type EditTaskDTO struct {
	Title        *string    `json:"title,omitempty" doc:"Task's title" maxLength:"50"`
	Description  *string    `json:"description,omitempty" doc:"Task's description"`
	DueDate      *time.Time `json:"due_date,omitempty" doc:"Task's due date" format:"date-time"`
	AssigneeName *string    `json:"assignee_name,omitempty" doc:"Name of the person to assign the task to"`
	Status       *string    `json:"status,omitempty" doc:"Task's status"`
}

// ChangeTaskStatusDTO represents the input for changing a task's status
type ChangeTaskStatusDTO struct {
	Status string `json:"status" doc:"New status for the task" minLength:"1"`
}
