# AI SDK


<p align="left">
    <a href="https://github.com/unsafe0x0/ai-sdk/releases/tag/v1.3.0">
        <img src="https://img.shields.io/badge/v1.3.0-blue.svg" alt="v1.3.0">
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
- Gemini (`GeminiProvider`) currently does not support options.

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
│  ├── options.go       # Options type for request customization
│  ├── provider.go       # Provider interface and SDK wrapper
│  └── response.go       # JSON response parser
providers/               # Provider implementations
│  ├── openrouter.go     # OpenRouter provider
│  ├── groqcloud.go      # GroqCloud provider
│  ├── mistral.go        # Mistral provider
│  ├── openai.go         # OpenAI provider
│  ├── perplexity.go     # Perplexity provider
│  ├── anthropic.go      # Anthropic provider
│  └── gemini.go         # Gemini provider
stream/                  # Streaming response parsing
│  └── stream.go         # Stream chunk parser
example/                 # Example programs
   └── readme.md         # Example usage of the SDK
```

## Declaring Providers

To use a provider, you first need to initialize it with your API key and a model name. Here are examples for each supported provider:

```go
// GroqCloud
client := ai.NewSDK(ai.NewGroqCloudProvider("YOUR_GROQ_API_KEY", "llama3-8b-8192"))

// Mistral
client := ai.NewSDK(ai.NewMistralProvider("YOUR_MISTRAL_API_KEY", "mistral-small-latest"))

// OpenRouter
client := ai.NewSDK(ai.NewOpenRouterProvider("YOUR_OPEN_ROUTER_API_KEY", "openrouter/sonoma-dusk-alpha"))

// OpenAI
client := ai.NewSDK(ai.NewOpenAiProvider("YOUR_OPENAI_API_KEY", "gpt-3.5-turbo"))

// Perplexity
client := ai.NewSDK(ai.NewPerplexityProvider("YOUR_PERPLEXITY_API_KEY", "sonar-pro"))

// Anthropic
client := ai.NewSDK(ai.NewAnthropicProvider("YOUR_ANTHROPIC_API_KEY", "claude-3.5"))

// Gemini
client := ai.NewSDK(ai.NewGeminiProvider("YOUR_GEMINI_API_KEY", "gemini-2.5-flash"))
```

## Streaming Completions

This SDK supports streaming responses from providers. You can use either `StreamComplete` for simple streaming or `StreamCompleteWithOptions` to include additional parameters.

### Basic Streaming

The `StreamComplete` method takes a context, a slice of messages, and a callback function to handle each chunk of the response.

```go
err := client.StreamComplete(ctx, messages, func(chunk string) error {
	fmt.Print(chunk)
	return nil
})
```

### Streaming with Options

The `StreamCompleteWithOptions` method allows you to pass additional options to the provider.

```go
var opts ai.Options
maxTokens := 2048
opts.MaxTokens = &maxTokens

err := client.StreamCompleteWithOptions(ctx, messages, func(chunk string) error {
	fmt.Print(chunk)
	return nil
}, &opts)
```

### Available Options

The `Options` struct supports the following fields:

- `MaxTokens` (\*int): The maximum number of tokens to generate.
- `ReasoningEffort` (\*int): A custom parameter to control reasoning effort (e.g., 1 for low, 2 for medium, 3 for high). Note that not all providers support this field.

## Non Streaming Completions

For scenarios where you need the full response at once, you can use the non streaming methods.

### Basic Completion

The `Complete` method blocks until the full response is received and returns it as a single string.

```go
resp, err := client.Complete(ctx, messages)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(resp)
```

### Completion with Options

Similarly, `CompleteWithOptions` allows you to specify options for a non-streaming request.

```go
var opts ai.Options
maxTokens := 1024
opts.MaxTokens = &maxTokens

resp, err := client.CompleteWithOptions(ctx, messages, &opts)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(resp)
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
