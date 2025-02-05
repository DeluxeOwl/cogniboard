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
