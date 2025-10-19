package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/unsafe0x0/ai/v2"
	"github.com/unsafe0x0/ai/v2/sdk"
)

func main() {
	_ = godotenv.Load()

	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	if geminiApiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}

	client := ai.Gemini(geminiApiKey)

	// Store conversation history
	messages := []sdk.Message{}

	// Create a scanner for reading user input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("AI Chat (type 'exit' to quit)")
	fmt.Println("=" + string(make([]byte, 40)) + "=")

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		if userInput == "exit" || userInput == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if userInput == "" {
			continue
		}

		// Add user message to history
		messages = append(messages, sdk.Message{
			Role:    "user",
			Content: userInput,
		})

		// Make API request with streaming
		res := client.ChatCompletion(context.Background(), &sdk.CompletionRequest{
			Messages:        messages,
			Model:           "gemini-2.0-flash-exp",
			SystemPrompt:    "You are a helpful assistant.",
			Temperature:     0.7,
			ReasoningEffort: "low",
			Stream:          false,
			Tools: map[string]sdk.Tool{
				"weather_tool": {
					Description: "Get the weather in a location",
					InputSchema: sdk.InputSchema{
						"location": {Type: "string", Description: "The location to get the weather for", Required: true},
					},
					Execute: func(ctx context.Context, args json.RawMessage) (any, error) {
						return 25, nil
					},
				},
				"code_executor": {
					Description: "Execute code in a given language and return output",
					InputSchema: sdk.InputSchema{
						"language": {Type: "string", Description: "Programming language (go, python, js)", Required: true},
						"code":     {Type: "string", Description: "The code snippet to execute", Required: true},
					},
					Execute: func(ctx context.Context, args json.RawMessage) (any, error) {
						var params struct {
							Language string `json:"language"`
							Code     string `json:"code"`
						}
						if err := json.Unmarshal(args, &params); err != nil {
							return "", err
						}

						var cmd *exec.Cmd
						switch params.Language {
						case "python":
							cmd = exec.CommandContext(ctx, "python3", "-c", params.Code)
						case "go":
							tmpFile := "temp_exec.go"
							if err := os.WriteFile(tmpFile, []byte(params.Code), 0644); err != nil {
								return "", err
							}
							cmd = exec.CommandContext(ctx, "go", "run", tmpFile)
						case "js", "javascript":
							cmd = exec.CommandContext(ctx, "node", "-e", params.Code)
						default:
							return "", fmt.Errorf("unsupported language: %s", params.Language)
						}

						output, err := cmd.CombinedOutput()
						if err != nil {
							return "", fmt.Errorf("execution failed: %v\nOutput: %s", err, string(output))
						}

						res, _ := json.Marshal(map[string]string{
							"language": params.Language,
							"output":   string(output),
						})
						return string(res), nil
					},
				},
				"fetch_web_page": sdk.Tool{
					Description: "Fetches a webpage using curl and returns its HTML content.",
					InputSchema: sdk.InputSchema{
						"url": {
							Type:        "string",
							Description: "the url of the webpage to fetch",
						},
					},
					Execute: func(ctx context.Context, args json.RawMessage) (any, error) {
						var params struct {
							URL string `json:"url"`
						}

						if err := json.Unmarshal(args, &params); err != nil {
							return "", err
						}

						cmd := exec.CommandContext(ctx, "curl", "-sL", fmt.Sprintf("https://r.jina.ai/%s", params.URL))
						out, err := cmd.CombinedOutput()
						if err != nil {
							return "", err
						}
						return string(out), nil
					},
				},
			},
		})

		// Check for errors
		if res.Error != nil {
			log.Printf("Error: %v\n", res.Error)
			// Remove the last message since it failed
			messages = messages[:len(messages)-1]
			continue
		}

		// Read and display streaming response
		fmt.Print("AI: ")
		var assistantResponse string

		if res.Stream != nil {
			defer res.Stream.Close()

			buf := make([]byte, 1024)
			for {
				n, err := res.Stream.Read(buf)
				if n > 0 {
					chunk := string(buf[:n])
					fmt.Print(chunk)
					assistantResponse += chunk
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("\nStream error: %v\n", err)
					break
				}
			}
			fmt.Println()
		} else {
			assistantResponse = res.Content
			fmt.Print(assistantResponse)
		}

		// Add assistant response to history
		if assistantResponse != "" {
			messages = append(messages, sdk.Message{
				Role:    "assistant",
				Content: assistantResponse,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v\n", err)
	}
}
