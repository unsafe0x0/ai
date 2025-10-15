package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/unsafe0x0/ai/base"
	"github.com/unsafe0x0/ai/sdk"
)

type MistralProvider struct {
	*base.Provider
	APIKey string
}

func NewMistralProvider(apiKey string) *MistralProvider {
	p := &MistralProvider{
		APIKey: apiKey,
	}
	p.Provider = &base.Provider{APICaller: p}
	return p
}

func (p *MistralProvider) CallAPI(ctx context.Context, messages []sdk.Message, streamMode bool, opts *sdk.Options) (io.ReadCloser, error) {
	url := "https://api.mistral.ai/v1/chat/completions"

	chatMessages := []map[string]string{}
	for _, m := range messages {
		chatMessages = append(chatMessages, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}

	body := map[string]interface{}{
		"messages": chatMessages,
		"stream":   streamMode,
	}
	if opts != nil {
		if opts.Model != "" {
			body["model"] = opts.Model
		}
		if opts.MaxCompletionTokens != 0 {
			body["max_tokens"] = opts.MaxCompletionTokens
		}
		if opts.Temperature != 0 {
			body["temperature"] = opts.Temperature
		}
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

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
		return nil, &sdk.APIError{
			StatusCode: resp.StatusCode,
			Message:    string(b),
			Body:       b,
		}
	}

	return resp.Body, nil
}

func (p *MistralProvider) ParseStream(body io.Reader, onChunk func(string) error) error {
	reader := bufio.NewReader(body)

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			if bytes.HasPrefix(line, []byte("data: ")) {
				line = line[len("data: "):]
			}
			if bytes.Equal(line, []byte("[DONE]")) {
				return nil
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal(line, &chunk); err == nil {
				for _, c := range chunk.Choices {
					if c.Delta.Content != "" {
						if err := onChunk(c.Delta.Content); err != nil {
							return err
						}
					}
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
