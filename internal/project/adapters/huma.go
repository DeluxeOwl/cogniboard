package adapters

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/DeluxeOwl/cogniboard/internal/project/app"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/commands"
	"github.com/danielgtaylor/huma/v2"
)

// Huma handles HTTP requests using the huma framework
type Huma struct {
	api         huma.API
	app         *app.Application
	fileStorage project.FileStorage
}

// NewHuma creates a new huma HTTP server
func NewHuma(api huma.API, app *app.Application, fileStorage project.FileStorage) *Huma {
	return &Huma{api: api, app: app, fileStorage: fileStorage}
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
},
) (*struct{}, error) {
	taskID, err := project.NewTaskID()
	if err != nil {
		return nil, err
	}

	data := input.RawBody.Data()

	cmd := commands.CreateTask{
		TaskID: taskID,
	}

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

	err = h.app.Commands.CreateTask.Handle(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("task create: %w", err)
	}

	var domainFiles []project.File
	now := time.Now()
	for _, files := range input.RawBody.Form.File {
		for _, file := range files {
			filename := file.Filename
			body, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("read file %s: %w", filename, err)
			}

			err = h.fileStorage.Store(ctx, taskID, filename, body)
			if err != nil {
				return nil, fmt.Errorf("store file %s: %w", filename, err)
			}

			domainFiles = append(domainFiles, project.File{
				Name:       filename,
				Size:       file.Size,
				MimeType:   file.Header.Get("Content-Type"),
				UploadedAt: now,
			})
		}
	}

	err = h.app.Commands.AttachFilesToTask.Handle(ctx, commands.AttachFilesToTask{
		TaskID: taskID,
		Files:  domainFiles,
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
		dtos[i] = project.ToInTaskDTO(&task)
	}

	return &struct{ Body project.InTasksDTO }{
		Body: project.InTasksDTO{Tasks: dtos},
	}, nil
}

func (h *Huma) editTask(ctx context.Context, input *struct {
	TaskID  string `path:"taskId"`
	RawBody huma.MultipartFormFiles[project.InCreateTaskDTO]
},
) (*struct{}, error) {
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
},
) (*struct{}, error) {
	err := h.app.Commands.ChangeTaskStatus.Handle(ctx, commands.ChangeTaskStatus{
		TaskID: input.TaskID,
		Status: input.Body.Status,
	})

	return nil, handleError(err)
}
