package adapters

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/DeluxeOwl/cogniboard/internal/project/app"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/commands"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/queries"
	"github.com/danielgtaylor/huma/v2"
)

// Huma handles HTTP requests using the huma framework
type Huma struct {
	api    huma.API
	app    *app.Application
	logger *slog.Logger
}

// NewHuma creates a new huma HTTP server
func NewHuma(api huma.API, app *app.Application, logger *slog.Logger) *Huma {
	return &Huma{api: api, app: app, logger: logger}
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
		OperationID: "project-chat",
		Method:      http.MethodPost,
		Path:        "/chat",
		Summary:     "Chat about your project",
	}, h.chatWithProject)

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

// In DTOs - for input adapters: e.g REST api

func (h *Huma) chatWithProject(ctx context.Context, input *struct{ Body queries.ChatWithProject }) (*huma.StreamResponse, error) {
	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {

			ctx.SetHeader("Content-Type", "text/my-stream")
			writer := ctx.BodyWriter()

			if d, ok := writer.(interface{ SetWriteDeadline(time.Time) error }); ok {
				d.SetWriteDeadline(time.Now().Add(5 * time.Second))
			}

			stream, err := h.app.Queries.ChatWithProject.Handle(ctx.Context(), queries.ChatWithProject{
				Messages: input.Body.Messages,
			})
			if err != nil {
				h.logger.Error("create stream", "err", err)
			}

			for stream.Next() {
				chunk := stream.Current()

				writer.Write(chunk)
				if f, ok := writer.(http.Flusher); ok {
					f.Flush()
				} else {
					fmt.Println("error: unable to flush")
				}
			}
		},
	}, nil
}

// note: huma doesnt play well with struct embedding
type CreateTask struct {
	Title       string          `form:"title" doc:"Task's name" minLength:"1" maxLength:"50" required:"true"`
	Description string          `form:"description" doc:"Task's description"`
	DueDate     time.Time       `form:"due_date" doc:"Task's due date (if any)" format:"date-time"`
	Assignee    string          `form:"assignee_name" doc:"Task's asignee (if any)"`
	Files       []huma.FormFile `form:"files"`
}

func (h *Huma) createTask(ctx context.Context, input *struct {
	RawBody huma.MultipartFormFiles[CreateTask]
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

	if data.Assignee != "" {
		cmd.AssigneeName = &data.Assignee
	}

	if data.Description != "" {
		cmd.Description = &data.Description
	}

	filesToUpload, err := h.prepareFilesForUpload(input.RawBody.Form.File)
	if err != nil {
		return nil, handleError(err)
	}

	err = h.app.Commands.CreateTask.Handle(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("task create: %w", err)
	}

	return nil, handleError(h.app.Commands.
		AttachFilesToTask.
		Handle(ctx, commands.AttachFilesToTask{
			TaskID: taskID,
			Files:  filesToUpload,
		}))
}

func (h *Huma) prepareFilesForUpload(rawFiles map[string][]*multipart.FileHeader) ([]commands.FileToUpload, error) {
	var filesToUpload []commands.FileToUpload
	for _, files := range rawFiles {
		for _, file := range files {
			filename := file.Filename

			fileReader, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("open file %s: %w", filename, err)
			}

			domainFile, err := project.NewFile(filename, file.Size)
			if err != nil {
				return nil, fmt.Errorf("create file %s: %w", filename, err)
			}

			filesToUpload = append(filesToUpload, commands.FileToUpload{
				Metadata: domainFile,
				Content:  fileReader,
			})
		}
	}
	return filesToUpload, nil
}

type ListTasks struct {
	Tasks []Task `json:"tasks"`
}
type Task struct {
	project.TaskSnapshot
	Files []File `json:"files"`
}

type File struct {
	project.FileSnapshot
}

func (h *Huma) getTasks(ctx context.Context, input *struct{}) (*struct{ Body ListTasks }, error) {
	tasks, err := h.app.Queries.AllTasks.Handle(ctx, struct{}{})
	if err != nil {
		return nil, huma.Error400BadRequest("couldn't get tasks", err)
	}

	dtos := make([]Task, len(tasks))
	for i, task := range tasks {
		dtos[i] = taskFrom(&task)
	}

	return &struct{ Body ListTasks }{
		Body: ListTasks{Tasks: dtos},
	}, nil
}

type EditTask struct {
	Title        string          `form:"title" doc:"Task's name" minLength:"1" maxLength:"50" required:"true"`
	Description  string          `form:"description" doc:"Task's description"`
	DueDate      time.Time       `form:"due_date" doc:"Task's due date (if any)" format:"date-time"`
	AssigneeName string          `form:"assignee_name" doc:"Task's asignee (if any)"`
	Files        []huma.FormFile `form:"files"`
}

func (h *Huma) editTask(ctx context.Context, input *struct {
	TaskID  string `path:"taskId"`
	RawBody huma.MultipartFormFiles[EditTask]
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

type ChangeTaskStatus struct {
	Status string `json:"status" doc:"New status for the task" minLength:"1"`
}

func (h *Huma) changeTaskStatus(ctx context.Context, input *struct {
	TaskID string           `path:"taskId"`
	Body   ChangeTaskStatus `json:"body"`
},
) (*struct{}, error) {
	err := h.app.Commands.ChangeTaskStatus.Handle(ctx, commands.ChangeTaskStatus{
		TaskID: input.TaskID,
		Status: input.Body.Status,
	})

	return nil, handleError(err)
}

func taskFrom(task *project.Task) Task {
	snap := task.GetSnapshot()
	dto := Task{
		TaskSnapshot: *task.GetSnapshot(),
		Files:        make([]File, len(snap.Files)),
	}

	for i := range snap.Files {
		dto.Files[i] = fileFrom(&snap.Files[i])
	}

	return dto
}

func fileFrom(file *project.File) File {
	return File{
		file.GetSnapshot(),
	}
}
