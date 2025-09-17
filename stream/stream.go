package stream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
)

// this is only used for Gemini streaming response parsing
type GeminiChunk struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func Stream(body io.Reader, onChunk func(string) error) error {
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

			// special handling for Gemini streaming response
			var chunk GeminiChunk
			if err := json.Unmarshal(line, &chunk); err == nil {
				if len(chunk.Candidates) > 0 && len(chunk.Candidates[0].Content.Parts) > 0 {
					transformedChunk := map[string]interface{}{
						"choices": []map[string]interface{}{
							{
								"delta": map[string]interface{}{
									"content": chunk.Candidates[0].Content.Parts[0].Text,
								},
							},
						},
					}
					transformedLine, _ := json.Marshal(transformedChunk)
					if err := onChunk(string(transformedLine)); err != nil {
						return err
					}
					continue
				}
			}

			if err := onChunk(string(line)); err != nil {
				return err
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
