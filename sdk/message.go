package sdk

type Message struct {
	Role       string            `json:"role"`
	Content    string            `json:"content"`
	ToolCallID string            `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCallRequest `json:"tool_calls,omitempty"`
}

type CompletionResponse struct {
	Content   string
	ToolCalls []ToolCallRequest
	Role      string
}
