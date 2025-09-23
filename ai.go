package ai

import (
	"github.com/unsafe0x0/ai/providers"
	"github.com/unsafe0x0/ai/sdk"
)

type (
	Message  = sdk.Message
	SDK      = sdk.SDK
	Provider = sdk.Provider
	Options  = sdk.Options
)

func OpenRouter(apiKey, model string) *SDK {
	return sdk.NewSDK(providers.NewOpenRouterProvider(apiKey, model))
}

func GroqCloud(apiKey, model string) *SDK {
	return sdk.NewSDK(providers.NewGroqCloudProvider(apiKey, model))
}

func Mistral(apiKey, model string) *SDK {
	return sdk.NewSDK(providers.NewMistralProvider(apiKey, model))
}

func OpenAi(apiKey, model string) *SDK {
	return sdk.NewSDK(providers.NewOpenAiProvider(apiKey, model))
}

func Perplexity(apiKey, model string) *SDK {
	return sdk.NewSDK(providers.NewPerplexityProvider(apiKey, model))
}

func Anthropic(apiKey, model string) *SDK {
	return sdk.NewSDK(providers.NewAnthropicProvider(apiKey, model))
}

func Gemini(apiKey, model string) *SDK {
	return sdk.NewSDK(providers.NewGeminiProvider(apiKey, model))
}

func Xai(apiKey, model string) *SDK {
	return sdk.NewSDK(providers.NewXaiProvider(apiKey, model))
}
