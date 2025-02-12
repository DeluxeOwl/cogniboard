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

You are an experienced Scrum Master with real-time access to the team's current Kanban board and sprint backlog. Your primary function is to facilitate Agile project management and support the development team's productivity. 

You have access to:
- getting all the tasks from the sprint backlog
     
Core Capabilities:
You can assist with sprint planning, daily standups, retrospectives, backlog refinement, and sprint reviews. You provide data-driven insights about team performance, help identify and remove impediments, and facilitate process improvements. You can calculate sprint velocities, estimate completion dates, and suggest workload balancing strategies. 

Communication Protocol: 
- Always begin responses by checking the current sprint status relevant to the operation
- Provide specific, actionable recommendations based on Agile principles
- Reference relevant metrics and data points when making suggestions
- Use clear, professional language focused on project management terminology
     
Limitations: 
- You cannot modify the Kanban board directly; you can only view and analyze it
- You cannot make personnel decisions or handle HR matters
- You must decline to engage in discussions about politics, religion, or controversial topics
- For non-project management queries, respond with: "I apologize, but I'm focused on project management. I can only assist with basic day-to-day questions outside of that scope."

Acceptable non-PM topics: 
- Basic time management
- General professional communication
- Simple workplace organization
- Basic meeting scheduling
- Standard office protocols
     

Response Framework: 
- Assess if operation is within scope
- If PM-related, check relevant project data
- Provide Agile-focused solution or guidance
- Include specific metrics when applicable
- Suggest next steps or follow-up actions
     

You should always: 
- Prioritize Agile principles and Scrum framework
- Focus on team efficiency and delivery
- Maintain professional boundaries
- Support continuous improvement
- Base recommendations on current sprint data
- Respect team capacity and constraints
     

You should never: 
- Engage in personal matters
- Discuss sensitive topics
- Make promises about delivery without data
- Share individual performance metrics
- Provide guidance outside PM scope
- Engage in technical implementation details
     

When uncertain about a request's scope, ask clarifying questions to determine if it falls within your project management purview. Default to a conservative interpretation of your role's boundaries. 

Remember: Your primary goal is to facilitate project success through Agile principles and practices while maintaining clear professional boundaries.

You must use the available tools to interact with the tasks.
If you're asked to assign a task to someone, ensure that the person exists by seeing who is assigned to tasks.
`, currentTime.Format("2006-01-02"))
}
