package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type Provider interface {
	CreateCompletion(ctx context.Context, messages []Message, opts *Options) (*CompletionResponse, error)
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
	Tools           map[string]Tool
	MaxToolSteps    int                                         // for preventing infinite loop error, defaults to 5
	onToolCall      func(toolName string, args json.RawMessage) // for ui callbacks
}

func (sdk *SDK) ChatCompletion(ctx context.Context, req *CompletionRequest) *Response {

	if req.MaxToolSteps == 0 {
		req.MaxToolSteps = 5
	}

	opts := &Options{
		Model:               req.Model,
		SystemPrompt:        req.SystemPrompt,
		MaxCompletionTokens: req.MaxTokens,
		Temperature:         req.Temperature,
		ReasoningEffort:     req.ReasoningEffort,
		Tools:               req.Tools,
		MaxToolSteps:        req.MaxToolSteps,
	}

	// for now tool calls only work for full responses

	if req.Stream {
		return sdk.streamingCompletion(ctx, req.Messages, opts)
	}

	if len(req.Tools) > 0 {
		return sdk.chatCompletionWithTools(ctx, req.Messages, opts, req.Tools, req.onToolCall)
	}

	if !req.Stream {
		return sdk.simpleCompletion(ctx, req.Messages, opts)
	}

	stream, err := sdk.provider.CreateCompletionStream(ctx, req.Messages, opts)
	if err != nil {
		return &Response{Error: err}
	}

	return &Response{Stream: &Stream{reader: stream}}
}

func (sdk *SDK) simpleCompletion(ctx context.Context, messages []Message, opts *Options) *Response {
	compResp, err := sdk.provider.CreateCompletion(ctx, messages, opts)
	return &Response{Content: compResp.Content, Error: err}
}

func (sdk *SDK) streamingCompletion(ctx context.Context, messages []Message, opts *Options) *Response {
	stream, err := sdk.provider.CreateCompletionStream(ctx, messages, opts)
	if err != nil {
		return &Response{Error: err}
	}
	return &Response{Stream: &Stream{reader: stream}}
}

func (sdk *SDK) chatCompletionWithTools(
	ctx context.Context,
	initialMessages []Message,
	opts *Options,
	tools map[string]Tool,
	onToolCall func(string, json.RawMessage),
) *Response {
	messages := append([]Message{}, initialMessages...)

	for step := 0; step < opts.MaxToolSteps; step++ {
		compResp, err := sdk.provider.CreateCompletion(ctx, messages, opts)

		if err != nil {
			return &Response{Error: err}
		}

		if len(compResp.ToolCalls) == 0 {
			return &Response{Content: compResp.Content}
		}

		messages = append(messages, Message{
			Role:      "assistant",
			Content:   compResp.Content,
			ToolCalls: compResp.ToolCalls,
		})

		for _, toolCall := range compResp.ToolCalls {
			tool, exists := tools[toolCall.Name]
			if !exists {
				messages = append(messages, Message{
					Role:       "tool",
					ToolCallID: toolCall.ID,
					Content:    fmt.Sprintf(`{"error": "Tool" '%s' not found"}`, toolCall.Name),
				})
				continue
			}

			if onToolCall != nil {
				onToolCall(toolCall.Name, toolCall.Arguments)
			}

			result, err := tool.Execute(ctx, toolCall.Arguments)

			var resultContent string
			if err != nil {
				resultContent = fmt.Sprintf(`{"error": "%s"}`, err.Error())
			} else {
				resultBytes, marshalErr := json.Marshal(result)
				if marshalErr != nil {
					resultContent = fmt.Sprintf(`{"error": "Failes to marshal toolCall result: %s"}`, marshalErr.Error())
				} else {
					resultContent = string(resultBytes)
				}
			}

			messages = append(messages, Message{
				Role:       "tool",
				ToolCallID: toolCall.ID,
				Content:    resultContent,
			})
		}
	}
	return &Response{
		Error: fmt.Errorf("reached maximum tool steps (%d) without final answer", opts.MaxToolSteps),
	}
}
