package sdk

import (
	"context"
	"errors"
)

type Provider interface {
	Generate(ctx context.Context, messages []Message, opts *Options, onChunk func(string) error) (string, error)
}

type SDK struct {
	provider Provider
}

func NewSDK(provider Provider) *SDK {
	return &SDK{provider: provider}
}

func (sdk *SDK) Generate(ctx context.Context, messages []Message, opts *Options) (string, error) {
	return sdk.provider.Generate(ctx, messages, opts, nil)
}

func (sdk *SDK) GenerateStream(ctx context.Context, messages []Message, opts *Options, onChunk func(string) error) (string, error) {
	if onChunk == nil {
		return "", errors.New("onChunk callback must not be nil")
	}
	return sdk.provider.Generate(ctx, messages, opts, onChunk)
}
