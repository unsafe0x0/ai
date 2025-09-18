package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/unsafe0x0/ai/base"
	"github.com/unsafe0x0/ai/sdk"
)

type AnthropicProvider struct {
	*base.Provider
	APIKey string
	Model  string
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	p := &AnthropicProvider{
		APIKey: apiKey,
		Model:  model,
	}
	p.Provider = &base.Provider{APICaller: p}
	return p
}

func (p *AnthropicProvider) CallAPI(ctx context.Context, messages []sdk.Message, streamMode bool, opts *sdk.Options) (io.ReadCloser, error) {
	url := "https://api.anthropic.com/v1/messages"

	chatMessages := []map[string]string{}
	for _, m := range messages {
		chatMessages = append(chatMessages, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}

	var systemPrompt string
	if len(chatMessages) > 0 && chatMessages[0]["role"] == "system" {
		systemPrompt = chatMessages[0]["content"]
		chatMessages = chatMessages[1:]
	}

	body := map[string]interface{}{
		"model":      p.Model,
		"system":     systemPrompt,
		"messages":   chatMessages,
		"stream":     streamMode,
		"max_tokens": 4096,
	}
	if opts != nil {
		if opts.MaxTokens != 0 {
			body["max_tokens"] = opts.MaxTokens
		}
		if opts.ReasoningEffort != "" {
			body["reasoning_effort"] = opts.ReasoningEffort
		}
		if opts.Temperature != 0 {
			body["temperature"] = opts.Temperature
		}
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", p.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("anthropic error: %s", string(b))
	}

	return resp.Body, nil
}
