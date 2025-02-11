package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DeluxeOwl/cogniboard/internal/project"
	"github.com/DeluxeOwl/cogniboard/internal/project/app/operations"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/shared"
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
func NewOpenAIAdapter(client *openai.Client, config OpenAIConfig) operations.ChatService {
	return &openAIAdapter{
		client: client,
		config: config,
	}
}

func convertMessages(messages []operations.Message) ([]openai.ChatCompletionMessageParamUnion, error) {
	result := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		var combinedText strings.Builder
		for _, content := range msg.Content {
			if content.Type == "text" {
				combinedText.WriteString(content.Text)
				combinedText.WriteString("\n")
			}
		}

		// Has to be this way for non-openai providers compatibility
		// Open AI does it like this now: {content:{type: text, text:"..."}}
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

func (a *openAIAdapter) StreamChat(ctx context.Context, messages []operations.Message, tools []project.ChatTool) (project.StreamingChunk, error) {
	converted, err := convertMessages(messages)
	if err != nil {
		return nil, fmt.Errorf("failed to convert messages: %w", err)
	}

	params := openai.ChatCompletionNewParams{
		Messages: openai.F(converted),
		Model:    openai.F(a.config.Model),
	}

	hasTools := len(tools) != 0
	if !hasTools {
		stream := a.client.Chat.Completions.NewStreaming(ctx, params)

		return &streamingChunk{
			stream: stream,
		}, nil
	}

	openAIToolParams := make([]openai.ChatCompletionToolParam, len(tools))
	for i, chatTool := range tools {
		converted := convertChatTool(chatTool)
		openAIToolParams[i] = *converted
	}
	params.Tools = openai.F(openAIToolParams)

	stream := a.client.Chat.Completions.NewStreaming(ctx, params)
	acc := openai.ChatCompletionAccumulator{}
	extractor := NewToolCallExtractor()
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		extractor.ExtractToolCallsFromChoices(chunk.Choices)

		if tool, ok := acc.JustFinishedToolCall(); ok {
			for _, providedTool := range tools {
				providedToolFuncName := providedTool.GetFuncName()

				if callID, ok := extractor.GetFunctionCallID(providedToolFuncName); ok {
					result, err := providedTool.CallHandler(ctx, tool.Arguments)
					if err != nil {
						return nil, err
					}

					if delta, ok := extractor.GetDelta(providedToolFuncName); ok {
						// Extract tool call into assistant message
						assistantMsg := extractor.ExtractToolCallIntoAssistantMessage(delta)
						params.Messages.Value = append(params.Messages.Value, assistantMsg)
						params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(callID, result))

					}
				}
			}

			// Start a new stream with the updated messages
			stream = a.client.Chat.Completions.NewStreaming(ctx, params)
			acc = openai.ChatCompletionAccumulator{}
		}
	}

	return nil, errors.New("couldnt match any stream")
}

func convertChatTool(tool project.ChatTool) *openai.ChatCompletionToolParam {
	funcParams := shared.FunctionParameters{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
	for _, param := range tool.GetToolParams() {
		props := funcParams["properties"].(map[string]any)
		props[param.Name] = map[string]string{
			"type": param.ParamType,
		}

		if param.Required {
			req := funcParams["required"].([]string)
			funcParams["required"] = append(req, param.Name)
		}

	}
	return &openai.ChatCompletionToolParam{
		Type: openai.F(openai.ChatCompletionToolType("object")),
		Function: openai.F(shared.FunctionDefinitionParam{
			Name:        openai.String(tool.GetFuncName()),
			Description: openai.String(tool.GetFuncDescription()),
			Parameters:  openai.F(funcParams),
		}),
	}
}

type toolCallExtractor struct {
	functionNamesToCallIDs map[string]string
	functionNamesToDelta   map[string]openai.ChatCompletionChunkChoicesDelta
}

func NewToolCallExtractor() *toolCallExtractor {
	return &toolCallExtractor{
		functionNamesToCallIDs: map[string]string{},
		functionNamesToDelta:   map[string]openai.ChatCompletionChunkChoicesDelta{},
	}
}

func (e *toolCallExtractor) GetFunctionCallID(functionName string) (string, bool) {
	id, ok := e.functionNamesToCallIDs[functionName]
	return id, ok
}

func (e *toolCallExtractor) GetDelta(functionName string) (openai.ChatCompletionChunkChoicesDelta, bool) {
	delta, ok := e.functionNamesToDelta[functionName]
	return delta, ok
}

func (e *toolCallExtractor) ExtractToolCallsFromChoices(choices []openai.ChatCompletionChunkChoice) {
	for _, choice := range choices {
		delta := choice.Delta
		for _, toolCall := range delta.ToolCalls {
			if toolCall.Type == openai.ChatCompletionChunkChoicesDeltaToolCallsTypeFunction {
				e.functionNamesToCallIDs[toolCall.Function.Name] = toolCall.ID
				e.functionNamesToDelta[toolCall.Function.Name] = delta
			}
		}
	}
}

// ExtractToolCallIntoAssistantMessage creates an assistant message from a tool call delta
func (e *toolCallExtractor) ExtractToolCallIntoAssistantMessage(delta openai.ChatCompletionChunkChoicesDelta) openai.ChatCompletionAssistantMessageParam {
	return openai.ChatCompletionAssistantMessageParam{
		Role: openai.F(openai.ChatCompletionAssistantMessageParamRoleAssistant),
		ToolCalls: openai.F([]openai.ChatCompletionMessageToolCallParam{
			{
				ID:   openai.String(delta.ToolCalls[0].ID),
				Type: openai.F(openai.ChatCompletionMessageToolCallType(delta.ToolCalls[0].Type)),
				Function: openai.F(openai.ChatCompletionMessageToolCallFunctionParam{
					Name:      openai.String(delta.ToolCalls[0].Function.Name),
					Arguments: openai.String(delta.ToolCalls[0].Function.Arguments),
				}),
			},
		}),
	}
}
