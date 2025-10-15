package sdk

type Options struct {
	Model               string  `json:"model,omitempty"`
	SystemPrompt        string  `json:"system_prompt,omitempty"`
	MaxCompletionTokens int     `json:"max_completion_tokens,omitempty"`
	ReasoningEffort     string  `json:"reasoning_effort,omitempty"`
	Temperature         float32 `json:"temperature,omitempty"`
}
