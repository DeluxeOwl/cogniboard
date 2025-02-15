package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/DeluxeOwl/cogniboard/internal/project"
)

type ChatWithProjectHandler decorator.OperationHandler[ChatWithProject, project.StreamingChunk]

type chatWithProjectHandler struct {
	chatService ChatService
	repo        project.TaskRepository
	embeddings  project.EmbeddingStorage
}

// ChatWithProjectReadModel defines the interface for reading chat interactions
type ChatWithProjectReadModel interface {
	ChatWithProject(ctx context.Context) (project.StreamingChunk, error)
}

func NewChatWithProjectHandler(
	chatService ChatService,
	repo project.TaskRepository,
	logger *slog.Logger,
	embeddings project.EmbeddingStorage,
) ChatWithProjectHandler {
	return decorator.ApplyOperationDecorators(
		&chatWithProjectHandler{chatService: chatService, repo: repo, embeddings: embeddings},
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
	StreamChat(
		ctx context.Context,
		messages []Message,
		tools []project.ChatTool,
	) (project.StreamingChunk, error)
}

type EditTaskArgs struct {
	TaskID       string     `json:"taskID"`
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	DueDate      *time.Time `json:"dueDate"`
	AssigneeName *string    `json:"assigneeName"`
	Status       *string    `json:"status"`
}

type SearchDocumentsForTaskArgs struct {
	TaskID string `json:"taskID"`
	Query  string `json:"query"`
}

type SearchAllDocumentsArgs struct {
	Query string `json:"query"`
}

// TODO: chat should be a domain object. For now it's fine to do some business logic here.
// The openai adapter should depend on the message domain entity
func (h *chatWithProjectHandler) Handle(
	ctx context.Context,
	operation ChatWithProject,
) (project.StreamingChunk, error) {
	operation.enrichWithSystemPrompt()

	return h.chatService.StreamChat(ctx, operation.Messages, []project.ChatTool{
		project.Tool[SearchAllDocumentsArgs]{
			FuncName:    "search_all_documents",
			Description: "Searches through all documents, attached to ANY task based on embeddings, gets back the most likely results for the user's query",
			Params: []project.ToolParam{
				{
					Name:      "query",
					ParamType: "string",
					Required:  true,
				},
			},
			Handler: func(ctx context.Context, sdfta SearchAllDocumentsArgs) (string, error) {
				res, err := h.embeddings.SearchAllDocuments(ctx, sdfta.Query)
				if err != nil {
					return "couldn't search documents", fmt.Errorf("search documents: %w", err)
				}

				marshaled, err := json.Marshal(res)
				if err != nil {
					return "couldn't search documents", fmt.Errorf("search documents: %w", err)
				}

				return string(marshaled), nil
			},
		},
		project.Tool[SearchDocumentsForTaskArgs]{
			FuncName:    "search_documents_for_task",
			Description: "Searches through all documents attached to a task based on embeddings, gets back the embedding search for the user's query",
			Params: []project.ToolParam{
				{
					Name:      "taskID",
					ParamType: "string",
					Required:  true,
				},
				{
					Name:      "query",
					ParamType: "string",
					Required:  true,
				},
			},
			Handler: func(ctx context.Context, sdfta SearchDocumentsForTaskArgs) (string, error) {
				res, err := h.embeddings.SearchDocumentsForTask(
					ctx,
					project.TaskID(sdfta.TaskID),
					sdfta.Query,
				)
				if err != nil {
					return "couldn't search the task", fmt.Errorf("search docs for task: %w", err)
				}

				marshaled, err := json.Marshal(res)
				if err != nil {
					return "couldn't search the task", fmt.Errorf("search docs for task: %w", err)
				}

				return string(marshaled), nil
			},
		},
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
				err := h.repo.UpdateTask(
					ctx,
					project.TaskID(cmd.TaskID),
					func(t *project.Task) (*project.Task, error) {
						var status *project.TaskStatus
						if cmd.Status != nil {
							s := project.TaskStatus(*cmd.Status)
							status = &s
						}

						if err := t.Edit(cmd.Title, cmd.Description, cmd.DueDate, cmd.AssigneeName, status); err != nil {
							return nil, err
						}
						return t, nil
					},
				)
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
	additionalInstruction := "You MUST not talk about anything else other than things related to the team's kanban board and the sprint's backlog, and other project related things (such as tech questions)."
	// additionalInstruction := ""
	return fmt.Sprintf(`
You are CogniMaster, an AI assistant designed to function as a Scrum Master with real-time access to a team's Kanban board and sprint backlog. Your primary role is to facilitate Agile project management and support the development team's productivity. 

%s

Current context:
<current_time>
%s
</current_time>

<existent_assignees>
%s
<existent_assignees>

Available tools:
1. get_tasks: Retrieve information for all tasks (ID, title, description, due date, assignee, timestamps, status, associated files).
2. edit_task: Edit a specific task.
3. search_documents_for_task: Searches through all documents attached to a task based on embeddings, gets back the embedding search for the user's query
4. search_all_documents: Searches through all documents, attached to ANY task based on embeddings, gets back the most likely results for the user's query

Instructions:
1. Analyze the user's message and determine the appropriate action.
2. If you need task information, use the get_tasks tool.
3. If you need to edit a task, use the edit_task tool. Always use get_tasks first if you don't have enough information to use edit_task.
4. After editing a task, end your response with "@refetch" to update the sprint board in real-time.
5. Provide clear, direct answers without announcing your thought process or using formulaic starts.
6. For complex problems, break them down systematically but present the solution conversationally.
7. If asked by the user "where can I find this information" - respond with the task with the attached files or where you got the information from in detail: the task id, its status, and to whom the task is assigned to
8. CogniMaster can only assign tasks to existent assignees

Before responding, organize your thoughts inside <analysis> tags to ensure a clear, non-repetitive response. Consider the following:
- Summarize the main request or question from the user
- List the information and tools needed to address this request
- Plan the structure of your response to be clear and concise
- Consider any potential challenges or edge cases

Now, please process the user's message and provide an appropriate response.
Below is the user's message:
`, additionalInstruction, currentTime.Format("2006-01-02"), strings.Join([]string{"John", "Mary", "Steve", "Laura", "Alex"}, ","))
}
