package sdk

import "context"

type Provider interface {
	Complete(ctx context.Context, messages []Message) (string, error)
	StreamComplete(ctx context.Context, messages []Message, onChunk func(string) error) error
	CompleteWithOptions(ctx context.Context, messages []Message, opts *Options) (string, error)
	StreamCompleteWithOptions(ctx context.Context, messages []Message, onChunk func(string) error, opts *Options) error
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

func (sdk *SDK) CompleteWithOptions(ctx context.Context, messages []Message, opts *Options) (string, error) {
	return sdk.provider.CompleteWithOptions(ctx, messages, opts)
}

func (sdk *SDK) StreamCompleteWithOptions(ctx context.Context, messages []Message, onChunk func(string) error, opts *Options) error {
	return sdk.provider.StreamCompleteWithOptions(ctx, messages, onChunk, opts)
}
