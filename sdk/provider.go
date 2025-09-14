package sdk

import "context"

type Provider interface {
	Complete(ctx context.Context, messages []Message) (string, error)
	StreamComplete(ctx context.Context, messages []Message, onChunk func(string) error) error
}

type SDK struct {
	provider Provider
}

func NewSDK(provider Provider) *SDK {
	return &SDK{provider: provider}
}

func (sdk *SDK) Complete(ctx context.Context, messages []Message) (string, error) {
	return sdk.provider.Complete(ctx, messages)
}

func (sdk *SDK) StreamComplete(ctx context.Context, messages []Message, onChunk func(string) error) error {
	return sdk.provider.StreamComplete(ctx, messages, onChunk)
}
