package adapters

import (
	"context"
	"net/http"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/DeluxeOwl/cogniboard/internal/project/app"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/commands"
	"github.com/danielgtaylor/huma/v2"
)

// Huma handles HTTP requests using the huma framework
type Huma struct {
	api huma.API
	app *app.Application
}

// NewHuma creates a new huma HTTP server
func NewHuma(api huma.API, app *app.Application) *Huma {
	return &Huma{api: api, app: app}
}

// Register registers all HTTP routes with huma
func (h *Huma) Register() {
	huma.Register(h.api, huma.Operation{
		OperationID: "task-create",
		Method:      http.MethodPost,
		Path:        "/tasks/create",
		Summary:     "Create a task",
	}, h.createTask)

	huma.Register(h.api, huma.Operation{
		OperationID: "tasks",
		Method:      http.MethodGet,
		Path:        "/tasks",
		Summary:     "Get all tasks",
	}, h.getTasks)

	huma.Register(h.api, huma.Operation{
		OperationID: "task-edit",
		Method:      http.MethodPost,
		Path:        "/tasks/{taskId}/edit",
		Summary:     "Edit a task",
	}, h.editTask)

	huma.Register(h.api, huma.Operation{
		OperationID: "task-change-status",
		Method:      http.MethodPost,
		Path:        "/tasks/{taskId}/status",
		Summary:     "Change task status",
	}, h.changeTaskStatus)
}

// handleError is a helper function to handle errors consistently
func handleError(err error) error {
	if err != nil {
		return huma.Error422UnprocessableEntity("operation failed", err)
	}
	return nil
}

func (h *Huma) createTask(ctx context.Context, input *struct{ Body project.InCreateTaskDTO }) (*struct{}, error) {
	err := h.app.Commands.CreateTask.Handle(ctx, commands.CreateTask{
		Title:        input.Body.Title,
		Description:  input.Body.Description,
		DueDate:      input.Body.DueDate,
		AssigneeName: input.Body.AssigneeName,
	})

	return nil, handleError(err)
}

func (h *Huma) getTasks(ctx context.Context, input *struct{}) (*struct{ Body project.InTasksDTO }, error) {
	tasks, err := h.app.Queries.AllTasks.Handle(ctx, struct{}{})
	if err != nil {
		return nil, huma.Error400BadRequest("couldn't get tasks", err)
	}

	dtos := make([]project.InTaskDTO, len(tasks))
	for i, task := range tasks {
		dtos[i] = project.InTaskDTO{
			ID:          string(task.ID()),
			Title:       task.Title(),
			Description: task.Description(),
			DueDate:     task.DueDate(),
			Assignee:    task.Asignee(),
			CreatedAt:   task.CreatedAt(),
			CompletedAt: task.CompletedAt(),
			UpdatedAt:   task.UpdatedAt(),
			Status:      string(task.Status()),
		}
	}

	return &struct{ Body project.InTasksDTO }{
		Body: project.InTasksDTO{Tasks: dtos},
	}, nil
}

func (h *Huma) editTask(ctx context.Context, input *struct {
	TaskID string                `path:"taskId"`
	Body   project.InEditTaskDTO `json:"body"`
}) (*struct{}, error) {
	err := h.app.Commands.EditTask.Handle(ctx, commands.EditTask{
		TaskID:       input.TaskID,
		Title:        input.Body.Title,
		Description:  input.Body.Description,
		DueDate:      input.Body.DueDate,
		AssigneeName: input.Body.AssigneeName,
		Status:       input.Body.Status,
	})

	return nil, handleError(err)
}

func (h *Huma) changeTaskStatus(ctx context.Context, input *struct {
	TaskID string                        `path:"taskId"`
	Body   project.InChangeTaskStatusDTO `json:"body"`
}) (*struct{}, error) {
	err := h.app.Commands.ChangeTaskStatus.Handle(ctx, commands.ChangeTaskStatus{
		TaskID: input.TaskID,
		Status: input.Body.Status,
	})

	return nil, handleError(err)
}
