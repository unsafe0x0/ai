package stream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
)

func StreamChunks(body io.Reader, onChunk func(string) error) error {
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

			var chunkObj map[string]any
			if err := json.Unmarshal(line, &chunkObj); err != nil {
				continue
			}
			if choices, ok := chunkObj["choices"].([]any); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]any); ok {
					if delta, ok := choice["delta"].(map[string]any); ok {
						if content, ok := delta["content"].(string); ok {
							if err := onChunk(content); err != nil {
								return err
							}
						}
					}
					if message, ok := choice["message"].(map[string]any); ok {
						if content, ok := message["content"].(string); ok {
							if err := onChunk(content); err != nil {
								return err
							}
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
