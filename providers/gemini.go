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

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

type GenerationConfig struct {
	Temperature     float32 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type GeminiRequest struct {
	Contents          []GeminiContent   `json:"contents"`
	SystemInstruction *GeminiContent    `json:"system_instruction,omitempty"`
	GenerationConfig  *GenerationConfig `json:"generation_config,omitempty"`
}

type Candidate struct {
	Content      GeminiContent `json:"content"`
	FinishReason string        `json:"finishReason"`
}

type GeminiResponseChunk struct {
	Candidates []Candidate `json:"candidates"`
}

type PromptFeedback struct {
	BlockReason string `json:"blockReason"`
}

type GeminiResponse struct {
	Candidates     []Candidate     `json:"candidates"`
	PromptFeedback *PromptFeedback `json:"promptFeedback,omitempty"`
}

type ContentBlockedError struct {
	Reason string
	Body   []byte
}

func (e *ContentBlockedError) Error() string {
	return fmt.Sprintf("content blocked by safety filters. Finish Reason: %s. Response Body: %s", e.Reason, string(e.Body))
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
		cfg.Temperature = 0.7

		if opts.Temperature > 0 {
			cfg.Temperature = opts.Temperature
		}

		if opts.MaxCompletionTokens > 0 {
			cfg.MaxOutputTokens = opts.MaxCompletionTokens
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

	if !streamMode {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var response GeminiResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to parse non-streaming JSON response: %w. Body: %s", err, string(body))
		}

		if response.PromptFeedback != nil && response.PromptFeedback.BlockReason != "" {
			return nil, &ContentBlockedError{
				Reason: response.PromptFeedback.BlockReason,
				Body:   body,
			}
		}

		if len(response.Candidates) > 0 {
			candidate := response.Candidates[0]

			if candidate.FinishReason == "SAFETY" || candidate.FinishReason == "RECITATION" {
				return nil, &ContentBlockedError{
					Reason: candidate.FinishReason,
					Body:   body,
				}
			}

			var fullText string
			for _, part := range candidate.Content.Parts {
				fullText += part.Text
			}

			if fullText == "" {

				return nil, fmt.Errorf("non-streaming response body was successfully parsed but contained empty text (FinishReason: %s). Raw body: %s", candidate.FinishReason, string(body))
			}

			return io.NopCloser(bytes.NewReader([]byte(fullText))), nil
		}

		return nil, fmt.Errorf("non-streaming response body was successfully parsed but contained no candidates. Raw body: %s", string(body))
	}

	return resp.Body, nil
}

func (p *GeminiProvider) ParseResponse(body io.Reader, onChunk func(string) error) error {
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

			var chunk GeminiResponseChunk

			if jsonErr := json.Unmarshal(line, &chunk); jsonErr == nil {
				if len(chunk.Candidates) > 0 {
					candidate := chunk.Candidates[0]

					if candidate.FinishReason == "SAFETY" || candidate.FinishReason == "RECITATION" {
						return &ContentBlockedError{Reason: candidate.FinishReason, Body: line}
					}

					if len(candidate.Content.Parts) > 0 {
						text := candidate.Content.Parts[0].Text
						if text != "" {
							if chunkErr := onChunk(text); chunkErr != nil {
								return chunkErr
							}
						}
					}

					if candidate.FinishReason == "STOP" || candidate.FinishReason == "MAX_TOKENS" {
						return nil
					}
				}
			} else {
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
