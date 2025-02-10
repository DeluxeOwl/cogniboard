package queries

import (
	"context"
	"log/slog"

	"github.com/DeluxeOwl/cogniboard/internal/decorator"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
)

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

type StreamingChunk interface {
	Close() error
	Current() []byte
	Err() error
	Next() bool
}

type streamingChunk struct {
	stream *ssestream.Stream[openai.ChatCompletionChunk]
}

func (sc *streamingChunk) Close() error {
	return sc.stream.Close()
}

func (sc *streamingChunk) Current() []byte {
	chunk := sc.stream.Current()
	json := chunk.JSON.RawJSON()

	rawByteChunk := make([]byte, 0, len(json)+1)
	rawByteChunk = append(rawByteChunk, json...)
	rawByteChunk = append(rawByteChunk, '\n')

	return rawByteChunk
}

func (sc *streamingChunk) Err() error {
	return sc.stream.Err()
}

func (sc *streamingChunk) Next() bool {
	return sc.stream.Next()
}

type ChatWithProjectHandler decorator.QueryHandler[ChatWithProject, StreamingChunk]

type chatWithProjectHandler struct {
	client *openai.Client
}

type ChatWithProjectReadModel interface {
	ChatWithProject(ctx context.Context) (StreamingChunk, error)
}

func NewChatWithProjectHandler(client *openai.Client, logger *slog.Logger) ChatWithProjectHandler {
	return decorator.ApplyQueryDecorators(
		&chatWithProjectHandler{client: client},
		logger,
	)
}

func convertMessages(rawMessages []Message) []openai.ChatCompletionMessageParamUnion {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(rawMessages))
	for i, msg := range rawMessages {
		isUser := msg.Role == "user"
		for _, content := range msg.Content {
			if isUser {
				messages[i] = openai.UserMessage(content.Text)
			} else {
				messages[i] = openai.AssistantMessage(content.Text)
			}
		}

	}
	return messages
}

func (h *chatWithProjectHandler) Handle(ctx context.Context, query ChatWithProject) (StreamingChunk, error) {
	stream := h.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(convertMessages(query.Messages)),
		Model:    openai.F(openai.ChatModelO3Mini),
	})

	return &streamingChunk{
		stream: stream,
	}, nil
}
