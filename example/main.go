package main

import (
	"ai/providers"
	"ai/sdk"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file loaded:", err)
	}

	apiKey := os.Getenv("OPEN_ROUTER_API_KEY")
	modelID := "openrouter/sonoma-sky-alpha" // Replace with your desired model ID
	systemInstruction := "You are unsafeai based on ai sdk built by unsafezero" // Optional: set a system instruction here

	var client *sdk.SDK
	openRouterProvider := providers.OpenRouterProvider{APIKey: apiKey, Model: modelID}
	client = sdk.NewSDK(&openRouterProvider)

	reader := bufio.NewReader(os.Stdin)
	var systemMessage sdk.Message
	if systemInstruction != "" {
		systemMessage = sdk.Message{Role: "system", Content: systemInstruction}
	}

	for {
		fmt.Print("\nPrompt: ")
		prompt, _ := reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)
		if prompt == "exit" {
			break
		}

		var messages []sdk.Message
		if systemMessage.Content != "" {
			messages = append(messages, systemMessage)
		}
		messages = append(messages, sdk.Message{Role: "user", Content: prompt})

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
