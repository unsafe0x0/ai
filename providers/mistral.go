package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/unsafe0x0/ai/sdk"
	"github.com/unsafe0x0/ai/stream"
)

type MistralProvider struct {
	APIKey string
	Model  string
}

func (p *MistralProvider) callAPI(ctx context.Context, messages []sdk.Message, streamMode bool) (io.ReadCloser, error) {
	url := "https://api.mistral.ai/v1/chat/completions"

	chatMessages := []map[string]string{}
	for _, m := range messages {
		chatMessages = append(chatMessages, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}

	body := map[string]interface{}{
		"model":    p.Model,
		"messages": chatMessages,
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

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("mistral error: %s", string(b))
	}

	return resp.Body, nil
}

func (p *MistralProvider) Complete(ctx context.Context, messages []sdk.Message) (string, error) {
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

func (p *MistralProvider) StreamComplete(ctx context.Context, messages []sdk.Message, onChunk func(string) error) error {
	body, err := p.callAPI(ctx, messages, true)
	if err != nil {
		return err
	}
	defer body.Close()

	return stream.StreamChunks(body, onChunk)
}
