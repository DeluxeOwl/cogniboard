package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type ChatWithProjectHandler decorator.OperationHandler[ChatWithProject, project.StreamingChunk]

type chatWithProjectHandler struct {
	chatService ChatService
	repo        project.TaskRepository
}

// ChatWithProjectReadModel defines the interface for reading chat interactions
type ChatWithProjectReadModel interface {
	ChatWithProject(ctx context.Context) (project.StreamingChunk, error)
}

func NewChatWithProjectHandler(chatService ChatService, repo project.TaskRepository, logger *slog.Logger) ChatWithProjectHandler {
	return decorator.ApplyOperationDecorators(
		&chatWithProjectHandler{chatService: chatService, repo: repo},
		logger,
	)
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ChatWithProject struct {
	Messages []Message `json:"messages"`
}

type ChatService interface {
	StreamChat(ctx context.Context, messages []Message, tools []project.ChatTool) (project.StreamingChunk, error)
}

type EditTaskArgs struct {
	TaskID       string     `json:"taskID"`
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	DueDate      *time.Time `json:"dueDate"`
	AssigneeName *string    `json:"assigneeName"`
	Status       *string    `json:"status"`
}

// TODO: chat should be a domain object. For now it's fine to do some business logic here.
// The openai adapter should depend on the message domain entity
func (h *chatWithProjectHandler) Handle(ctx context.Context, operation ChatWithProject) (project.StreamingChunk, error) {
	operation.enrichWithSystemPrompt()

	return h.chatService.StreamChat(ctx, operation.Messages, []project.ChatTool{
		project.Tool[struct{}]{
			FuncName:    "get_tasks",
			Description: "Get all available tasks from the board",
			Params:      []project.ToolParam{},
			Handler: func(ctx context.Context, _ struct{}) (string, error) {
				tasks, err := h.repo.AllTasks(ctx)
				if err != nil {
					return "", fmt.Errorf("get tasks: %w", err)
				}

				snaps := make([]*project.TaskSnapshot, len(tasks))
				for i, task := range tasks {
					snaps[i] = task.GetSnapshot()
				}

				marshaled, err := json.Marshal(snaps)
				if err != nil {
					return "", fmt.Errorf("marshal tasks snapshot: %w", err)
				}

				return string(marshaled), nil
			},
		},
		project.Tool[EditTaskArgs]{
			FuncName:    "edit_task",
			Description: "Edit a task, only the taskID is required. The allowed values for status are pending, in_progress and in_review. This function cannot be used to mark a task as complete.",
			Params: []project.ToolParam{
				{
					Name:      "taskID",
					ParamType: "string",
					Required:  true,
				},
				{
					Name:      "title",
					ParamType: "string",
				},
				{
					Name:      "description",
					ParamType: "string",
				},
				{
					Name:      "dueDate",
					ParamType: "string",
				},
				{
					Name:      "assigneeName",
					ParamType: "string",
				},
				{
					Name:      "status",
					ParamType: "string",
				},
			},
			Handler: func(ctx context.Context, cmd EditTaskArgs) (string, error) {
				err := h.repo.UpdateTask(ctx, project.TaskID(cmd.TaskID), func(t *project.Task) (*project.Task, error) {
					var status *project.TaskStatus
					if cmd.Status != nil {
						s := project.TaskStatus(*cmd.Status)
						status = &s
					}

					if err := t.Edit(cmd.Title, cmd.Description, cmd.DueDate, cmd.AssigneeName, status); err != nil {
						return nil, err
					}
					return t, nil
				})

				if err != nil {
					return "couldn't edit task", nil
				}

				return "Edited task", nil
			},
		},
	})
}

func (op *ChatWithProject) enrichWithSystemPrompt() {
	op.Messages = append([]Message{
		{
			Role: "system",
			Content: []Content{
				{
					Type: "text",
					Text: NewSystemPrompt(time.Now()),
				},
			},
		},
	}, op.Messages...)
}

func NewSystemPrompt(currentTime time.Time) string {
	return fmt.Sprintf(`
<context>
The current time is %s
</context>

You are an experienced CogniMaster, an AI assitant designed as a Scrum Master with real-time access to the team's current Kanban board and sprint backlog. Your primary function is to facilitate Agile project management and support the development team's productivity. 

CogniMaster uses the tools to interact with the sprint backlog
- CogniMaster MUST use the get_tasks to get the status of the sprint
- CogniMaster MUST use the get_tasks to see the available assignees
- CogniMaster MUST use the edit_task to edit a task
- CogniMaster MUST use get_tasks if he doesn't have enough information to use edit_task

<important-instructions>
After editing a task, CogniMaster MUST end his response with "@refetch" so that the sprint board is updated in real-time.
<important-instructions>

`, currentTime.Format("2006-01-02"))
}
