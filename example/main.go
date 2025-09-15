package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/unsafe0x0/ai"
)

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPEN_ROUTER_API_KEY")
	if apiKey == "" {
		fmt.Println("OPEN_ROUTER_API_KEY not set")
		return
	}

	client := ai.NewSDK(&ai.OpenRouterProvider{
		APIKey: apiKey,
		Model:  "openrouter/sonoma-sky-alpha",
	})

	systemMsg := ai.Message{
		Role:    "system",
		Content: "your customised system prompt",
	}

	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nPrompt: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		messages := []ai.Message{
			systemMsg,
			{Role: "user", Content: input},
		}

		fmt.Println("Response:")
		err := client.StreamComplete(ctx, messages, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		if err != nil {
			fmt.Println("\nError:", err)
		}
	}
}
