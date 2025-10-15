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

type GeminiProvider struct {
	*base.Provider
	APIKey string
}

func NewGeminiProvider(apiKey string) *GeminiProvider {
	p := &GeminiProvider{
		APIKey: apiKey,
	}
	p.Provider = &base.Provider{APICaller: p}
	return p
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	Temperature     float32  `json:"temperature,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	CandidateCount  int      `json:"candidateCount,omitempty"`
	StopSequences   []string `json:"stopSequences,omitempty"`
}

type GeminiRequest struct {
	Contents          []GeminiContent   `json:"contents"`
	SystemInstruction *GeminiContent    `json:"system_instruction,omitempty"`
	GenerationConfig  *GenerationConfig `json:"generationConfig,omitempty"`
}

func (p *GeminiProvider) CallAPI(ctx context.Context, messages []sdk.Message, streamMode bool, opts *sdk.Options) (io.ReadCloser, error) {

	var model string
	if opts != nil && opts.Model != "" {
		model = opts.Model
	}

	var url string
	if streamMode {
		url = fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:streamGenerateContent?alt=sse&key=%s", model, p.APIKey)
	} else {
		url = fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, p.APIKey)
	}

	var systemInstruction *GeminiContent
	var geminiContents []GeminiContent

	for _, msg := range messages {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		if role == "system" {
			systemInstruction = &GeminiContent{
				Parts: []GeminiPart{{Text: msg.Content}},
			}
		} else {
			geminiContents = append(geminiContents, GeminiContent{
				Role:  role,
				Parts: []GeminiPart{{Text: msg.Content}},
			})
		}
	}

	reqBody := GeminiRequest{
		Contents:          geminiContents,
		SystemInstruction: systemInstruction,
	}

	if opts != nil {
		cfg := &GenerationConfig{}
		if opts.MaxCompletionTokens > 0 {
			cfg.MaxOutputTokens = opts.MaxCompletionTokens
		}
		if opts.Temperature > 0 {
			cfg.Temperature = opts.Temperature
		} else {
			cfg.Temperature = 0.7
		}
		reqBody.GenerationConfig = cfg
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
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

func (p *GeminiProvider) ParseStream(body io.Reader, onChunk func(string) error) error {
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
				Candidates []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				} `json:"candidates"`
			}

			if err := json.Unmarshal(line, &chunk); err == nil {
				if len(chunk.Candidates) > 0 && len(chunk.Candidates[0].Content.Parts) > 0 {
					if err := onChunk(chunk.Candidates[0].Content.Parts[0].Text); err != nil {
						return err
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
