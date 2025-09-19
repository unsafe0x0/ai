package sdk

import "context"

type Provider interface {
	Generate(ctx context.Context, messages []Message, opts *Options, onChunk func(string) error) (string, error)
}

type SDK struct {
	provider Provider
}

func NewSDK(provider Provider) *SDK {
	return &SDK{provider: provider}
}

func (sdk *SDK) Generate(ctx context.Context, messages []Message, opts *Options, onChunk func(string) error) (string, error) {
	return sdk.provider.Generate(ctx, messages, opts, onChunk)
}
