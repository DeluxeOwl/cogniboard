package adapters

import (
	"context"
	"net/http"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/project/app"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/commands"
	"github.com/danielgtaylor/huma/v2"
)

type Huma struct {
	api huma.API
	app *app.Application
}

func NewHuma(api huma.API, app *app.Application) *Huma {
	return &Huma{api: api, app: app}
}

func (h *Huma) Register() {
	huma.Register(h.api, huma.Operation{
		OperationID: "create-task",
		Method:      http.MethodPost,
		Path:        "/tasks/create",
		Summary:     "Create a task",
	}, h.createTask)

	huma.Register(h.api, huma.Operation{
		OperationID: "get-tasks",
		Method:      http.MethodGet,
		Path:        "/tasks",
		Summary:     "Get all tasks",
	}, h.getTasks)

	huma.Register(h.api, huma.Operation{
		OperationID: "assign-task",
		Method:      http.MethodPost,
		Path:        "/tasks/{taskId}/assign",
		Summary:     "Assign a task to someone",
	}, h.assignTask)

	huma.Register(h.api, huma.Operation{
		OperationID: "unassign-task",
		Method:      http.MethodPost,
		Path:        "/tasks/{taskId}/unassign",
		Summary:     "Unassign a task",
	}, h.unassignTask)

	huma.Register(h.api, huma.Operation{
		OperationID: "change-task-status",
		Method:      http.MethodPost,
		Path:        "/tasks/{taskId}/status",
		Summary:     "Change task status",
	}, h.changeTaskStatus)
}

type CreateTaskDTO struct {
	Title        string     `json:"title" doc:"Task's name" minLength:"1" maxLength:"50"`
	Description  *string    `json:"description,omitempty" doc:"Task's description"`
	DueDate      *time.Time `json:"due_date,omitempty" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName *string    `json:"assignee_name,omitempty" doc:"Task's asignee (if any)"`
}

func (h *Huma) createTask(ctx context.Context, input *struct{ Body CreateTaskDTO }) (*struct{}, error) {
	err := h.app.Commands.CreateTask.Handle(ctx, commands.CreateTask{
		Title:        input.Body.Title,
		Description:  input.Body.Description,
		DueDate:      input.Body.DueDate,
		AssigneeName: input.Body.AssigneeName,
	})

	if err != nil {
		return nil, huma.Error422UnprocessableEntity("couldn't create task", err)
	}

	return nil, nil
}

type GetTasksDTO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Assignee    *string    `json:"assignee"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Status      string     `json:"status"`
}
type Tasks struct {
	Tasks []GetTasksDTO `json:"tasks"`
}

func (h *Huma) getTasks(ctx context.Context, input *struct{}) (*struct{ Body Tasks }, error) {
	tasks, err := h.app.Queries.AllTasks.Handle(ctx, struct{}{})
	if err != nil {
		return nil, huma.Error400BadRequest("couldn't get tasks", err)
	}

	dtos := make([]GetTasksDTO, len(tasks))
	for i, task := range tasks {
		dtos[i] = GetTasksDTO{
			ID:          string(task.ID()),
			Title:       task.Title(),
			Description: task.Description(),
			DueDate:     task.DueDate(),
			Assignee:    task.Asignee(),
			CreatedAt:   task.CreatedAt(),
			CompletedAt: task.CompletedAt(),
			Status:      string(task.Status()),
		}
	}

	return &struct{ Body Tasks }{
		Body: Tasks{dtos},
	}, nil
}

type AssignTaskDTO struct {
	AssigneeName string `json:"assignee_name" doc:"Name of the person to assign the task to" minLength:"1"`
}

func (h *Huma) assignTask(ctx context.Context, input *struct {
	TaskID string        `path:"taskId"`
	Body   AssignTaskDTO `json:"body"`
}) (*struct{}, error) {
	err := h.app.Commands.AssignTask.Handle(ctx, commands.AssignTask{
		TaskID:       input.TaskID,
		AssigneeName: input.Body.AssigneeName,
	})

	if err != nil {
		return nil, huma.Error422UnprocessableEntity("couldn't assign task", err)
	}

	return nil, nil
}

func (h *Huma) unassignTask(ctx context.Context, input *struct {
	TaskID string `path:"taskId"`
}) (*struct{}, error) {
	err := h.app.Commands.UnassignTask.Handle(ctx, commands.UnassignTask{
		TaskID: input.TaskID,
	})

	if err != nil {
		return nil, huma.Error422UnprocessableEntity("couldn't unassign task", err)
	}

	return nil, nil
}

type ChangeTaskStatusDTO struct {
	Status string `json:"status" doc:"New status for the task" minLength:"1"`
}

func (h *Huma) changeTaskStatus(ctx context.Context, input *struct {
	TaskID string              `path:"taskId"`
	Body   ChangeTaskStatusDTO `json:"body"`
}) (*struct{}, error) {
	err := h.app.Commands.ChangeTaskStatus.Handle(ctx, commands.ChangeTaskStatus{
		TaskID: input.TaskID,
		Status: input.Body.Status,
	})

	if err != nil {
		return nil, huma.Error422UnprocessableEntity("couldn't change task status", err)
	}

	return nil, nil
}
