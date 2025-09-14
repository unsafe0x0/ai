# AI SDK

A simple Go SDK interacting with LLM providers. Supports streaming completions, custom instructions, and easy provider integration.

## Features

- Streamed chat completions from multiple providers
- Easily switch between providers and models
- Set custom system instructions for each session

## Project Structure

```text
go.mod, go.sum         # Go module files
readme.md              # Project documentation
sdk/                   # Core SDK interfaces and types
│   ├── message.go     # Message type and roles
│   └── provider.go    # Provider interface and SDK wrapper
providers/             # Provider implementations
│   ├── openrouter.go
│   └── ...
stream/                # Streaming response parsing
│   └── stream.go
example/               # example programs
│   └── main.go        # Example usage of the SDK
```

## Adding Providers or Features

- Implement the `Provider` interface in `sdk/sdk.go` for new providers.
- Add new models or API keys in `.env` and update `main.go` as needed.
- To support more advanced chat flows, extend the `Message` struct and message array logic.

---

**Note:** This project is in early development. Features, and structure may change frequently.
