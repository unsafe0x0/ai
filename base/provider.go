package base

import (
	"context"
	"io"

	"github.com/unsafe0x0/ai/sdk"
	"github.com/unsafe0x0/ai/stream"
)

type APICaller interface {
	CallAPI(ctx context.Context, messages []sdk.Message, streamMode bool, opts *sdk.Options) (io.ReadCloser, error)
}

type Provider struct {
	APICaller
}

func (p *Provider) Generate(ctx context.Context, messages []sdk.Message, opts *sdk.Options, onChunk func(string) error) (string, error) {
	streamMode := onChunk != nil

	body, err := p.CallAPI(ctx, messages, streamMode, opts)
	if err != nil {
		return "", err
	}
	defer body.Close()

	if streamMode {
		return "", stream.Stream(body, func(s string) error {
			var response struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := sdk.ParseResponse([]byte(s), &response); err != nil {
				return nil
			}
			if len(response.Choices) > 0 {
				return onChunk(response.Choices[0].Delta.Content)
			}
			return nil
		})
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
	if err := sdk.ParseResponse(respBytes, &response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", nil
	}
	return response.Choices[0].Message.Content, nil
}
