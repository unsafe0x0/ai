package sdk

type Options struct {
	MaxCompletionTokens int     `json:"max_completion_tokens,omitempty"`
	ReasoningEffort     string  `json:"reasoning_effort,omitempty"`
	Temperature         float32 `json:"temperature,omitempty"`
	Stream              bool    `json:"stream,omitempty"`
}
