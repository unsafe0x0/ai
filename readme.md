# AI SDK v1.0.0

A simple Go SDK interacting with LLM providers. Supports streaming completions, custom instructions, and easy provider integration.

## Features

- Streamed chat completions from multiple providers
- Easily switch between providers and models
- Set custom system instructions for each session

## Project Structure

```text
go.mod, go.sum           # Go module files
LICENSE                  # License file
readme.md                # Project documentation

ai/                      # Main package entrypoint
│  └── ai.go             # Re-exports SDK and providers for single import
sdk/                     # Core SDK interfaces and types
│  ├── message.go        # Message type and roles
│  └── provider.go       # Provider interface and SDK wrapper
providers/               # Provider implementations
│  └── openrouter.go     # OpenRouter provider
stream/                  # Streaming response parsing
│  └── stream.go         # Stream chunk parser
example/                 # Example programs
   └── main.go           # Example usage of the SDK
```

## Adding Providers or Features

- Implement the `Provider` interface in `sdk/sdk.go` for new providers.
- Add new models or API keys in `.env` and update `main.go` as needed.
- To support more advanced chat flows, extend the `Message` struct and message array logic.

## How to Use

1. **Install as a dependency:**
   Add this SDK to your Go project using:

   ```sh
   go get github.com/unsafe0x0/ai
   ```

2. **Set up API keys:**
   Create a `.env` file in your project root:

   ```
   OPEN_ROUTER_API_KEY=your_openrouter_key
   ```

3. **Basic usage example:**

   ```go
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
   ```

4. **Streaming responses:**

   ```go
   err := client.StreamComplete(ctx, messages, func(chunk string) error {
   fmt.Print(chunk)
   return nil
   })

   ```

5. **Non-streaming responses:**

   ```go
   resp, err := client.Complete(ctx, messages)
   if err != nil {
   fmt.Println("Error:", err)
   return
   }
   fmt.Println(resp)

   ```

6. **Switching providers:**
   Implement the `Provider` interface for new providers, or use the built-in ones in `providers/`.

7. **Custom instructions:**
   Add system or user messages to the `messages` slice to control the conversation context.

---

**Note:** This project is in early development. Features, and structure may change frequently.
