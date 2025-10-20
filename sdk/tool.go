// converts common tool definitions for all providers
// each provider needs to implement its own structures

package sdk

import (
	"context"
	"encoding/json"
)

type Tool struct {
	Description string          `json:"description,omitempty"`
	InputSchema InputSchema     `json:"inputSchema,omitempty"`
	Execute     ToolExecuteFunc `json:"-,omitempty"`
}

type InputSchema map[string]Property

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

type ToolExecuteFunc func(ctx context.Context, args json.RawMessage) (any, error)

type ToolCall struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
	IsError    bool   `json:"is_error,omitempty"`
}

type ToolCallRequest struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}
