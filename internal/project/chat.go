package project

import (
	"context"
	"encoding/json"
	"fmt"
)

type StreamingChunk interface {
	Close() error
	Current() []byte
	Err() error
	Next() bool
}

type ChatTool interface {
	isTool()
	CallHandler(context.Context, string) (string, error)
	GetFuncName() string
	GetFuncDescription() string
	GetToolParams() []ToolParam
}
type Tool[Params any] struct {
	FuncName    string
	Description string
	Params      []ToolParam
	Handler     func(context.Context, Params) (string, error)
}

func (t Tool[Params]) isTool() {}
func (t Tool[Params]) CallHandler(ctx context.Context, providedJSON string) (string, error) {
	var p Params

	err := json.Unmarshal([]byte(providedJSON), &p)
	if err != nil {
		return "", fmt.Errorf("call handler: %w", err)
	}

	return t.Handler(ctx, p)
}
func (t Tool[Params]) GetFuncName() string        { return t.FuncName }
func (t Tool[Params]) GetFuncDescription() string { return t.Description }
func (t Tool[Params]) GetToolParams() []ToolParam { return t.Params }

type ToolParam struct {
	Name      string
	ParamType string
	Required  bool
}
