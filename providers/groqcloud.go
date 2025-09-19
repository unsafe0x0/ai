package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/unsafe0x0/ai/base"
	"github.com/unsafe0x0/ai/sdk"
)

type GroqCloudProvider struct {
	*base.Provider
	APIKey string
	Model  string
}

func NewGroqCloudProvider(apiKey, model string) *GroqCloudProvider {
	p := &GroqCloudProvider{
		APIKey: apiKey,
		Model:  model,
	}
	p.Provider = &base.Provider{APICaller: p}
	return p
}

func (p *GroqCloudProvider) CallAPI(ctx context.Context, messages []sdk.Message, streamMode bool, opts *sdk.Options) (io.ReadCloser, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

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
	if opts != nil {
		if opts.MaxCompletionTokens != 0 {
			body["max_completion_tokens"] = opts.MaxCompletionTokens
		}
		if opts.ReasoningEffort != "" {
			body["reasoning_effort"] = opts.ReasoningEffort
		}
		if opts.Temperature != 0 {
			body["temperature"] = opts.Temperature
		}
		if opts.Stream {
			body["stream"] = opts.Stream
		}
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
		return nil, fmt.Errorf("groq error: %s", string(b))
	}

	return resp.Body, nil
}
func (p *GroqCloudProvider) ParseStream(body io.Reader, onChunk func(string) error) error {
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
					FinishReason *string `json:"finish_reason"`
				} `json:"choices"`
			}

			if err := json.Unmarshal(line, &chunk); err == nil {
				if len(chunk.Choices) > 0 {
					d := chunk.Choices[0].Delta.Content
					if d != "" {
						if err := onChunk(d); err != nil {
							return err
						}
					}
					if chunk.Choices[0].FinishReason != nil &&
						*chunk.Choices[0].FinishReason == "stop" {
						return nil
					}
				}
				continue
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
