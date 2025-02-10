package openaiproxy

import (
	"encoding/json"
	"log/slog"
)

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Stream   bool      `json:"stream"`
	Messages []Message `json:"messages"`
}

// ChatBodyProcessor implements BodyProcessor for chat requests
type ChatBodyProcessor struct {
	logger *slog.Logger
}

func NewChatBodyProcessor(logger *slog.Logger) *ChatBodyProcessor {
	return &ChatBodyProcessor{
		logger: logger,
	}
}

func (p *ChatBodyProcessor) Process(body []byte) ([]byte, error) {
	p.logger.Info("processing request body", "body", string(body))

	var chatReq ChatRequest
	if err := json.Unmarshal(body, &chatReq); err != nil {
		p.logger.Error("failed to unmarshal chat request", "err", err)
		return body, err
	}

	p.logger.Info("processed chat request", "request", chatReq)

	// You can modify chatReq here if needed

	return body, nil
}
