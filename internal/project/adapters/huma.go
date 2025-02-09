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

	var maxBodyBytes int64 = 50 * 1024 * 1024

	huma.Register(h.api, huma.Operation{
		OperationID:  "task-create",
		Method:       http.MethodPost,
		Path:         "/tasks/create",
		Summary:      "Create a task",
		MaxBodyBytes: maxBodyBytes,
	}, h.createTask)

	huma.Register(h.api, huma.Operation{
		OperationID: "tasks",
		Method:      http.MethodGet,
		Path:        "/tasks",
		Summary:     "Get all tasks",
	}, h.getTasks)

	huma.Register(h.api, huma.Operation{
		OperationID:  "task-edit",
		Method:       http.MethodPost,
		Path:         "/tasks/{taskId}/edit",
		Summary:      "Edit a task",
		MaxBodyBytes: maxBodyBytes,
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

func (h *Huma) createTask(ctx context.Context, input *struct {
	RawBody huma.MultipartFormFiles[project.InCreateTaskDTO]
}) (*struct{}, error) {
	data := input.RawBody.Data()

	cmd := commands.CreateTask{}

	if data.Title != "" {
		cmd.Title = data.Title
	}

	if !data.DueDate.IsZero() {
		cmd.DueDate = &data.DueDate
	}

	if data.AssigneeName != "" {
		cmd.AssigneeName = &data.AssigneeName
	}

	if data.Description != "" {
		cmd.Description = &data.Description
	}

	err := h.app.Commands.CreateTask.Handle(ctx, cmd)

	return nil, handleError(err)
}

func (h *Huma) getTasks(ctx context.Context, input *struct{}) (*struct{ Body project.InTasksDTO }, error) {
	tasks, err := h.app.Queries.AllTasks.Handle(ctx, struct{}{})
	if err != nil {
		return nil, huma.Error400BadRequest("couldn't get tasks", err)
	}

	dtos := make([]project.InTaskDTO, len(tasks))
	for i, task := range tasks {
		dtos[i] = project.ToInTaskDTO(&task)
	}

	return &struct{ Body project.InTasksDTO }{
		Body: project.InTasksDTO{Tasks: dtos},
	}, nil
}

func (h *Huma) editTask(ctx context.Context, input *struct {
	TaskID  string `path:"taskId"`
	RawBody huma.MultipartFormFiles[project.InCreateTaskDTO]
}) (*struct{}, error) {
	data := input.RawBody.Data()

	cmd := commands.EditTask{
		TaskID: input.TaskID,
		Title:  &data.Title, // validated from form
	}

	if !data.DueDate.IsZero() {
		cmd.DueDate = &data.DueDate
	}

	if data.AssigneeName != "" {
		cmd.AssigneeName = &data.AssigneeName
	}

	if data.Description != "" {
		cmd.Description = &data.Description
	}

	err := h.app.Commands.EditTask.Handle(ctx, cmd)

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
