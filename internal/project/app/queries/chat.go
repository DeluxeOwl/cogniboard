package queries

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
)

// Message represents a chat message
type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

// Content represents message content
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// StreamingChunk represents a stream of chat completion chunks
type StreamingChunk interface {
	Close() error
	Current() []byte
	Err() error
	Next() bool
}

// ChatService defines the interface for chat operations
type ChatService interface {
	StreamChat(ctx context.Context, messages []Message) (StreamingChunk, error)
}

// ChatWithProject represents a chat query with project context
type ChatWithProject struct {
	Messages []Message `json:"messages"`
}

// ChatWithProjectHandler handles chat queries
type ChatWithProjectHandler decorator.QueryHandler[ChatWithProject, StreamingChunk]

type chatWithProjectHandler struct {
	chatService ChatService
}

// ChatWithProjectReadModel defines the interface for reading chat interactions
type ChatWithProjectReadModel interface {
	ChatWithProject(ctx context.Context) (StreamingChunk, error)
}

// NewChatWithProjectHandler creates a new chat query handler
func NewChatWithProjectHandler(chatService ChatService, logger *slog.Logger) ChatWithProjectHandler {
	return decorator.ApplyQueryDecorators(
		&chatWithProjectHandler{chatService: chatService},
		logger,
	)
}

// TODO: chat should be a domain object. For now it's fine to do some business logic here.
// The openai adapter should depend on the message domain entity
func (h *chatWithProjectHandler) Handle(ctx context.Context, query ChatWithProject) (StreamingChunk, error) {
	query.Messages = append([]Message{
		{
			Role: "system",
			Content: []Content{
				{
					Type: "text",
					Text: NewSystemPromptV2(time.Now()),
				},
			},
		},
	}, query.Messages...)

	return h.chatService.StreamChat(ctx, query.Messages)
}

func NewSystemPrompt(currentTime time.Time) string {
	return fmt.Sprintf(`
<context>
The current time is %s
</context>

You are an experienced Scrum Master with real-time access to the team's current Kanban board and sprint backlog. Your primary function is to facilitate Agile project management and support the development team's productivity. 

You have access to:
- currently you don't have access to the board
     
Core Capabilities:
You can assist with sprint planning, daily standups, retrospectives, backlog refinement, and sprint reviews. You provide data-driven insights about team performance, help identify and remove impediments, and facilitate process improvements. You can calculate sprint velocities, estimate completion dates, and suggest workload balancing strategies. 

Communication Protocol: 
- Always begin responses by checking the current sprint status relevant to the query
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
- Assess if query is within scope
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
`, currentTime.String())
}
func NewSystemPromptV2(currentTime time.Time) string {
	return fmt.Sprintf(`
<context>
The current time is %s
</context>
You're a helpful assistant
`, currentTime.String())
}
