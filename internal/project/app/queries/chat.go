package queries

import (
	"context"
	"log/slog"

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

func (h *chatWithProjectHandler) Handle(ctx context.Context, query ChatWithProject) (StreamingChunk, error) {
	return h.chatService.StreamChat(ctx, query.Messages)
}
