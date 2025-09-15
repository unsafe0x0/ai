package ai

import (
	"github.com/unsafe0x0/ai/providers"
	"github.com/unsafe0x0/ai/sdk"
)

type (
	Message  = sdk.Message
	SDK      = sdk.SDK
	Provider = sdk.Provider
)

var (
	NewSDK = sdk.NewSDK
)

type (
	OpenRouterProvider = providers.OpenRouterProvider
	GroqCloudProvider  = providers.GroqCloudProvider
	MistralProvider    = providers.MistralProvider
)
