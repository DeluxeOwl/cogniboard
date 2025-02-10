package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeluxeOwl/cogniboard/internal/project/app/queries"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
)

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// OpenAIConfig holds configuration for the OpenAI client
type OpenAIConfig struct {
	Model string
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

type openAIAdapter struct {
	client *openai.Client
	config OpenAIConfig
}

// NewOpenAIAdapter creates a new OpenAI adapter that implements ChatService
func NewOpenAIAdapter(client *openai.Client, config OpenAIConfig) queries.ChatService {
	return &openAIAdapter{
		client: client,
		config: config,
	}
}

func convertMessages(messages []queries.Message) ([]openai.ChatCompletionMessageParamUnion, error) {
	result := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		var combinedText strings.Builder
		for _, content := range msg.Content {
			if content.Type == "text" {
				combinedText.WriteString(content.Text)
				combinedText.WriteString("\n")
			}
		}

		switch msg.Role {
		case RoleUser:
			result = append(result, openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleUser),
				Content: openai.F(any(combinedText.String())),
			})
		case RoleAssistant:
			result = append(result, openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleAssistant),
				Content: openai.F(any(combinedText.String())),
			})
		case RoleSystem:
			result = append(result, openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleSystem),
				Content: openai.F(any(combinedText.String())),
			})
		default:
			return nil, fmt.Errorf("unknown role: %s", msg.Role)
		}
	}

	return result, nil
}

func (a *openAIAdapter) StreamChat(ctx context.Context, messages []queries.Message) (queries.StreamingChunk, error) {
	converted, err := convertMessages(messages)
	if err != nil {
		return nil, fmt.Errorf("failed to convert messages: %w", err)
	}

	stream := a.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(converted),
		Model:    openai.F(a.config.Model),
	})

	return &streamingChunk{
		stream: stream,
	}, nil
}
