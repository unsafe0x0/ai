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

var (
	NewSDK                = sdk.NewSDK
	NewOpenRouterProvider = providers.NewOpenRouterProvider
	NewGroqCloudProvider  = providers.NewGroqCloudProvider
	NewMistralProvider    = providers.NewMistralProvider
	NewOpenAiProvider     = providers.NewOpenAiProvider
	NewPerplexityProvider = providers.NewPerplexityProvider
	NewAnthropicProvider  = providers.NewAnthropicProvider
	NewGeminiProvider     = providers.NewGeminiProvider
)

type (
	OpenRouterProvider = providers.OpenRouterProvider
	GroqCloudProvider  = providers.GroqCloudProvider
	MistralProvider    = providers.MistralProvider
	OpenAiProvider     = providers.OpenAiProvider
	PerplexityProvider = providers.PerplexityProvider
	AnthropicProvider  = providers.AnthropicProvider
	GeminiProvider     = providers.GeminiProvider
)
