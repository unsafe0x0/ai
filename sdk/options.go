package sdk

// options holds optional parameters for provider calls.
type Options struct {
	MaxTokens       *int `json:"max_tokens,omitempty"`
	ReasoningEffort *int `json:"reasoning_effort,omitempty"`
	// add more optional fields as needed
}
