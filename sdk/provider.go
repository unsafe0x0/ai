package sdk

import (
	"context"
	"io"
)

type Provider interface {
	CreateCompletion(ctx context.Context, messages []Message, opts *Options) (string, error)
	CreateCompletionStream(ctx context.Context, messages []Message, opts *Options) (io.ReadCloser, error)
}

type SDK struct {
	provider Provider
}

func NewSDK(provider Provider) *SDK {
	return &SDK{
		provider: provider,
	}
}

type Response struct {
	Content string
	Stream  *Stream
	Error   error
}

type Stream struct {
	reader io.ReadCloser
}

func (s *Stream) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func (s *Stream) Close() error {
	return s.reader.Close()
}

type CompletionRequest struct {
	Messages        []Message
	Model           string
	SystemPrompt    string
	MaxTokens       int
	Temperature     float32
	ReasoningEffort string
	Stream          bool
}

func (sdk *SDK) ChatCompletion(ctx context.Context, req *CompletionRequest) *Response {
	opts := &Options{
		Model:               req.Model,
		SystemPrompt:        req.SystemPrompt,
		MaxCompletionTokens: req.MaxTokens,
		Temperature:         req.Temperature,
		ReasoningEffort:     req.ReasoningEffort,
	}

	if !req.Stream {
		content, err := sdk.provider.CreateCompletion(ctx, req.Messages, opts)
		return &Response{Content: content, Error: err}
	}

	stream, err := sdk.provider.CreateCompletionStream(ctx, req.Messages, opts)
	if err != nil {
		return &Response{Error: err}
	}

	return &Response{Stream: &Stream{reader: stream}}
}
