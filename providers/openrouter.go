package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/unsafe0x0/ai/sdk"
	"github.com/unsafe0x0/ai/stream"
)

type OpenRouterProvider struct {
	APIKey string
	Model  string
}

func (p *OpenRouterProvider) callAPI(ctx context.Context, messages []sdk.Message, streamMode bool) (io.ReadCloser, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	body := map[string]interface{}{
		"model":    p.Model,
		"messages": messages,
		"stream":   streamMode,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (p *OpenRouterProvider) Complete(ctx context.Context, messages []sdk.Message) (string, error) {
	body, err := p.callAPI(ctx, messages, false)
	if err != nil {
		return "", err
	}
	defer body.Close()

	respBytes, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(respBytes), nil
}

func (p *OpenRouterProvider) StreamComplete(ctx context.Context, messages []sdk.Message, onChunk func(string) error) error {
	body, err := p.callAPI(ctx, messages, true)
	if err != nil {
		return err
	}
	defer body.Close()

	return stream.StreamChunks(body, onChunk)
}
