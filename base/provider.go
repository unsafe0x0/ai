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

func (p *Provider) Complete(ctx context.Context, messages []sdk.Message) (string, error) {
	return p.CompleteWithOptions(ctx, messages, nil)
}

func (p *Provider) StreamComplete(ctx context.Context, messages []sdk.Message, onChunk func(string) error) error {
	return p.StreamCompleteWithOptions(ctx, messages, onChunk, nil)
}

func (p *Provider) CompleteWithOptions(ctx context.Context, messages []sdk.Message, opts *sdk.Options) (string, error) {
	body, err := p.CallAPI(ctx, messages, false, opts)
	if err != nil {
		return "", err
	}
	defer body.Close()

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

	return response.Choices[0].Message.Content, nil
}

func (p *Provider) StreamCompleteWithOptions(ctx context.Context, messages []sdk.Message, onChunk func(string) error, opts *sdk.Options) error {
	body, err := p.CallAPI(ctx, messages, true, opts)
	if err != nil {
		return err
	}
	defer body.Close()

	return stream.Stream(body, func(s string) error {
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
