package sdk

// options holds optional parameters for provider calls.
type Options struct {
	MaxTokens       int     `json:"max_tokens,omitempty"`
	ReasoningEffort string  `json:"reasoning_effort,omitempty"`
	Temperature     float32 `json:"temperature,omitempty"`
	Stream          bool    `json:"stream,omitempty"`
	// add more optional fields as needed
}
