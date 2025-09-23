# AI SDK

<p align="left">
    <a href="https://github.com/unsafe0x0/ai-sdk/releases/tag/v1.3.3">
        <img src="https://img.shields.io/badge/v1.3.3-blue.svg" alt="v1.3.3">
    </a>
    <img src="https://img.shields.io/badge/Go-00ADD8?logo=go&labelColor=white" alt="Go">
    <br/>
    <img src="https://img.shields.io/badge/GroqCloud-FF6F00" alt="GroqCloud">
    <img src="https://img.shields.io/badge/Mistral-1976D2" alt="Mistral">
    <img src="https://img.shields.io/badge/OpenRouter-43A047" alt="OpenRouter">
    <img src="https://img.shields.io/badge/OpenAI-6E4AFF" alt="OpenAI">
    <img src="https://img.shields.io/badge/Perplexity-00B8D4" alt="Perplexity">
    <img src="https://img.shields.io/badge/Anthropic-FF4081" alt="Anthropic">
    <img src="https://img.shields.io/badge/Gemini-7C4DFF" alt="Gemini">
    <img src="https://img.shields.io/badge/Xai-FFFFFF" alt="Xai">
</p>

A simple Go SDK for interacting with LLM providers. Supports streaming completions, custom instructions, and easy provider integration.

## Features

- Streamed chat completions from multiple providers
- Easily switch between providers and models
- Set custom system instructions

## Providers

- GroqCloud (`GroqCloudProvider`)
- Mistral (`MistralProvider`)
- OpenRouter (`OpenRouterProvider`)
- OpenAI (`OpenAiProvider`)
- Perplexity (`PerplexityProvider`)
- Anthropic (`AnthropicProvider`)
- Gemini (`GeminiProvider`) currently does not support options other than streaming.
- Xai (`XaiProvider`)

## Project Structure

```text
go.mod, go.sum           # Go module files
LICENSE                  # License file
readme.md                # Project documentation
ai.go                    # Main package entrypoint

base/
│  └── provider.go       # Base provider with shared logic
sdk/                     # Core SDK interfaces and types
│  ├── message.go        # Message type and roles
│  ├── options.go        # Options type for request customization
│  └── provider.go       # Provider interface and SDK wrapper
providers/               # Provider implementations
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

## Declaring Providers

To use a provider, initialize it with your API key and a model name using the provided constructor functions:

```go
// OpenRouter
client := ai.OpenRouter("YOUR_OPEN_ROUTER_API_KEY", "openrouter/sonoma-dusk-alpha")

// GroqCloud
client := ai.GroqCloud("YOUR_GROQ_API_KEY", "openai/gpt-oss-20b")

// Mistral
client := ai.Mistral("YOUR_MISTRAL_API_KEY", "mistral-small-latest")

// OpenAI
client := ai.OpenAi("YOUR_OPENAI_API_KEY", "gpt-3.5-turbo")

// Perplexity
client := ai.Perplexity("YOUR_PERPLEXITY_API_KEY", "sonar-pro")

// Anthropic
client := ai.Anthropic("YOUR_ANTHROPIC_API_KEY", "claude-3.5")

// Gemini
client := ai.Gemini("YOUR_GEMINI_API_KEY", "gemini-2.5-flash")

// Xai
client := ai.Xai("YOUR_XAI_API_KEY", "xai-1.5-base")
```

## Available Options

The `Options` struct supports the following fields (see `example/main.go` for usage):

- `MaxCompletionTokens` (int): The maximum number of tokens to generate. Optional; set to 0 to skip.
- `ReasoningEffort` (string): Custom reasoning effort (e.g., "low", "medium", "high"). Optional; set to empty string to skip. Not all providers support this field.
- `Temperature` (float32): Controls randomness of the output (0.0 to 1.0). Optional; set to 0 to skip.
- `Stream` (bool): Whether to stream the response. Optional; default is false.

### Declaring Options

To declare options, create an instance of the `Options` struct and set the desired fields. You can conditionally set fields based on your needs:

```go
var opts ai.Options
if maxTokens > 0 {
    opts.MaxCompletionTokens = maxTokens
}
if temp > 0 {
    opts.Temperature = temp
}
opts.Stream = true
```

## Streaming Completions

This SDK supports streaming responses from providers. You can use the Generate method with a callback function to handle streamed chunks:

```go
client.Generate(ctx, messages, &opts, func(chunk string) error {
    fmt.Print(chunk)
    return nil
})
```

## Examples

All code examples for this SDK can be found in the [ai-sdk-examples](https://github.com/unsafe0x0/ai-sdk-examples) repository.

## Contributing

Contributions are welcome!

### Pull Requests

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature-name`).
3.  Commit your changes (`git commit -m 'Add some feature'`).
4.  Push to the branch (`git push origin feature/your-feature-name`).
5.  Open a pull request.

### Issues

If you find a bug or have a feature request, please open an issue on GitHub.

---

**Note:** This project is in early development. Features, and structure may change frequently.
