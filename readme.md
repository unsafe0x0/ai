# AI SDK

<p align="left">
    <a href="https://github.com/unsafe0x0/ai-sdk/releases/tag/v2.0.0">
        <img src="https://img.shields.io/badge/v2.0.0-blue.svg" alt="v2.0.0">
    </a>
    <img src="https://img.shields.io/badge/Go-00ADD8?logo=go&labelColor=white" alt="Go">
    <img src="https://img.shields.io/badge/License-MIT-green" alt="License">
    <a href="https://github.com/unsafe0x0/ai">
        <img src="https://img.shields.io/github/stars/unsafe0x0/ai?style=social" alt="GitHub stars">
        <img src="https://img.shields.io/github/forks/unsafe0x0/ai?style=social" alt="GitHub forks">
        <img src="https://img.shields.io/github/issues/unsafe0x0/ai" alt="GitHub issues">
        <img src="https://img.shields.io/github/last-commit/unsafe0x0/ai" alt="Last commit">
    </a>
    <br/>
</p>

A simple Go SDK for interacting with LLM providers. Supports streaming completions, custom instructions, and easy provider integration.

## Features

- Chat completions (non-streaming and streaming)
- Easily switch between providers and models
- Options for customizing requests (model, system prompt, max tokens, temperature, reasoning effort)

## Providers

   <div align="left">
    <img src="https://img.shields.io/badge/GroqCloud-FF6F00" alt="GroqCloud">
    <img src="https://img.shields.io/badge/Mistral-1976D2" alt="Mistral">
    <img src="https://img.shields.io/badge/OpenRouter-43A047" alt="OpenRouter">
    <img src="https://img.shields.io/badge/OpenAI-6E4AFF" alt="OpenAI">
    <img src="https://img.shields.io/badge/Perplexity-00B8D4" alt="Perplexity">
    <img src="https://img.shields.io/badge/Anthropic-FF4081" alt="Anthropic">
    <img src="https://img.shields.io/badge/Gemini-7C4DFF" alt="Gemini">
    <img src="https://img.shields.io/badge/Xai-FFFFFF" alt="Xai">
    <img src="https://img.shields.io/badge/Anannas-FF6F00" alt="Anannas">
    <br/>
    </div>

---

- GroqCloud (`GroqCloud`)
- Mistral (`Mistral`)
- OpenRouter (`OpenRouter`)
- OpenAI (`OpenAi`)
- Perplexity (`Perplexity`)
- Anthropic (`Anthropic`)
- Gemini (`Gemini`)
- Xai (`Xai`)
- Anannas (`Anannas`)

## Project Structure

```text
go.mod                   # Go module file
LICENSE                  # License file
readme.md                # Project documentation
ai.go                    # Main package entrypoint

base/
│  └── provider.go       # Base provider with shared logic
sdk/                     # Core SDK interfaces and types
│  ├── errors.go         # API errors handling
│  ├── message.go        # Message type and roles
│  ├── options.go        # Options type for request customization
│  └── provider.go       # Provider interface and SDK wrapper
providers/               # Provider implementations
│  ├── anannas.go        # Anannas provider
│  ├── anthropic.go      # Anthropic provider
│  ├── gemini.go         # Gemini provider
│  ├── groqcloud.go      # GroqCloud provider
│  ├── mistral.go        # Mistral provider
│  ├── openai.go         # OpenAI provider
│  ├── openrouter.go     # OpenRouter provider
│  └── perplexity.go     # Perplexity provider
│  └── xai.go            # Xai provider
example/                 # Example usage of the SDK
│  └── readme.md
```

## Available Options

The `Options` struct supports the following fields (see `example/readme.md` for usage):

- `Model` (string): The model to use (e.g., "gpt-4o", "gemini-1.5").
- `SystemPrompt` (string): Custom system prompt to guide the AI's behavior.
- `MaxCompletionTokens` (int): The maximum number of tokens to generate.
- `ReasoningEffort` (string): Custom reasoning effort (e.g., "low", "medium", "high").
- `Temperature` (float32): Controls randomness of the output (0.0 to 1.0).

## Getting Started

Import the SDK in your Go project:

```go
import "github.com/unsafe0x0/ai"
```

## Declaring Providers

To use a provider, initialize it with your API key using the provided constructor functions:

```go
// OpenRouter
client := ai.OpenRouter("YOUR_OPEN_ROUTER_API_KEY")

// GroqCloud
client := ai.GroqCloud("YOUR_GROQ_API_KEY")

// Mistral
client := ai.Mistral("YOUR_MISTRAL_API_KEY")

// OpenAI
client := ai.OpenAi("YOUR_OPENAI_API_KEY")

// Perplexity
client := ai.Perplexity("YOUR_PERPLEXITY_API_KEY")

// Anthropic
client := ai.Anthropic("YOUR_ANTHROPIC_API_KEY")

// Gemini
client := ai.Gemini("YOUR_GEMINI_API_KEY")

// Xai
client := ai.Xai("YOUR_XAI_API_KEY")

// Anannas
client := ai.Anannas("YOUR_ANANNAS_API_KEY")
```

### Messages and Options

Messages are simple role/content pairs:

```go
messages := []ai.Message{
    {Role: "user", Content: "Hello, how are you?"},
}
```

### Declaring Options

To declare options, create an instance of the `Options` struct and set the desired fields. You can conditionally set fields based on your needs:

```go
model := "gpt-4o" // use your desired model
systemPrompt := "You are a helpful assistant." // use your custom system prompt
maxTokens := 1000 // use according to your needs
reasoningEffort := "medium" // use "low", "medium", or "high"
temp := 0.7 // use a temperature between 0.0 and 1.0

var opts ai.Options
if model != "" {
    opts.Model = model
}
if systemPrompt != "" {
        opts.SystemPrompt = systemPrompt
}
if maxTokens > 0 {
    opts.MaxCompletionTokens = maxTokens
}
if reasoningEffort != "" {
    opts.ReasoningEffort = reasoningEffort
}
if temp > 0 {
    opts.Temperature = temp
}
```

## Streaming Completions

This SDK supports both non-streaming and streaming responses.

Non-streaming:

```go
resp, err := client.Generate(ctx, messages, &opts)
if err != nil {
    // handle error
}
fmt.Println(resp)
```

Streaming (callback must not be nil):

```go
_, err := client.GenerateStream(ctx, messages, &opts, func(chunk string) error {
    fmt.Print(chunk)
    return nil
})
if err != nil {
    // handle error
}
```

## Examples

All code examples for this SDK can be found in the [ai-sdk-examples](https://github.com/unsafe0x0/ai-sdk-examples) repository.

## Contributing

Contributions are welcome!

### Pull Requests

1.  Fork the repository.
2.  Create a new branch.
3.  Commit your changes.
4.  Push to the branch.
5.  Open a pull request.

### Issues

If you find a bug or have a feature request, please open an issue on GitHub.

---

**Note:** This project is in early development. Features, and structure may change frequently.
