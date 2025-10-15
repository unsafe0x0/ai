package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/unsafe0x0/ai/sdk"
)

type APICaller interface {
	CallAPI(ctx context.Context, messages []sdk.Message, streamMode bool, opts *sdk.Options) (io.ReadCloser, error)
}

type StreamParser interface {
	ParseStream(body io.Reader, onChunk func(string) error) error
}

type Provider struct {
	APICaller
}

func (p *Provider) Generate(
	ctx context.Context,
	messages []sdk.Message,
	opts *sdk.Options,
	onChunk func(string) error,
) (string, error) {

	if opts != nil && opts.SystemPrompt != "" {
		if len(messages) == 0 || messages[0].Role != "system" {
			messages = append([]sdk.Message{{Role: "system", Content: opts.SystemPrompt}}, messages...)
		}
	}

	streamMode := onChunk != nil

	body, err := p.CallAPI(ctx, messages, streamMode, opts)
	if err != nil {
		return "", err
	}
	defer body.Close()

	if streamMode {
		if parser, ok := p.APICaller.(StreamParser); ok {
			var out strings.Builder
			err := parser.ParseStream(body, func(chunk string) error {
				out.WriteString(chunk)
				if onChunk != nil {
					return onChunk(chunk)
				}
				return nil
			})
			return out.String(), err
		}
		return "", fmt.Errorf("streaming not supported by this provider")
	}

	respBytes, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", nil
	}
	return response.Choices[0].Message.Content, nil
}
