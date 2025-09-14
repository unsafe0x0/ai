package providers

import (
	"ai/sdk"
	"ai/stream"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type OpenRouterProvider struct {
	APIKey string
	Model  string
}

func (p *OpenRouterProvider) Complete(ctx context.Context, messages []sdk.Message) (string, error) {
	var buf bytes.Buffer
	if err := p.StreamComplete(ctx, messages, func(chunk string) error {
		_, werr := buf.WriteString(chunk)
		return werr
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (p *OpenRouterProvider) StreamComplete(ctx context.Context, messages []sdk.Message, onChunk func(string) error) error {
	url := "https://openrouter.ai/api/v1/chat/completions"
	body := map[string]interface{}{
		"model":    p.Model,
		"messages": messages,
		"stream":   true,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return stream.StreamChunks(resp.Body, onChunk)
}
