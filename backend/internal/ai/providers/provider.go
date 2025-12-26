package providers

import (
	"context"
	"fmt"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`    // "user", "assistant", "system"
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Content      string `json:"content"`
	TokensInput  int    `json:"tokens_input"`
	TokensOutput int    `json:"tokens_output"`
	Model        string `json:"model"`
	FinishReason string `json:"finish_reason"`
}

// StreamChunk for streaming responses
type StreamChunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   error  `json:"error,omitempty"`
}

// Provider interface for AI chat providers
type Provider interface {
	// Name returns the provider name (e.g., "openai", "claude")
	Name() string

	// Chat performs a synchronous chat completion
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// ChatStream performs a streaming chat completion
	ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error)

	// ValidateAPIKey checks if the API key is valid
	ValidateAPIKey() error

	// AvailableModels returns list of available models
	AvailableModels() []string

	// GetModel returns the currently configured model
	GetModel() string
}

// ProviderConfig holds configuration for a provider
type ProviderConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
	BaseURL     string // Optional custom base URL
}

// DefaultConfig returns a ProviderConfig with sensible defaults
func DefaultConfig() ProviderConfig {
	return ProviderConfig{
		MaxTokens:   2048,
		Temperature: 0.7,
	}
}

// ProviderFactory creates and manages AI providers
type ProviderFactory struct {
	configs map[string]ProviderConfig
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		configs: make(map[string]ProviderConfig),
	}
}

// SetConfig sets the configuration for a provider
func (f *ProviderFactory) SetConfig(name string, config ProviderConfig) {
	f.configs[name] = config
}

// GetConfig gets the configuration for a provider
func (f *ProviderFactory) GetConfig(name string) (ProviderConfig, bool) {
	config, ok := f.configs[name]
	return config, ok
}

// GetProvider returns a configured provider instance
func (f *ProviderFactory) GetProvider(name string) (Provider, error) {
	config, ok := f.configs[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not configured", name)
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("API key not set for provider %s", name)
	}

	switch name {
	case "openai":
		return NewOpenAIProvider(config), nil
	case "claude":
		return NewClaudeProvider(config), nil
	case "groq":
		return NewGroqProvider(config), nil
	case "gemini":
		return NewGeminiProvider(config), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

// AvailableProviders returns list of supported provider names
func AvailableProviders() []string {
	return []string{"openai", "claude", "groq", "gemini"}
}

// ProviderInfo provides information about a provider
type ProviderInfo struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Models      []string `json:"models"`
	Configured  bool     `json:"configured"`
}

// GetProviderInfo returns information about providers
func (f *ProviderFactory) GetProviderInfo() []ProviderInfo {
	infos := []ProviderInfo{
		{
			Name:        "openai",
			DisplayName: "OpenAI",
			Models:      []string{"gpt-4-turbo", "gpt-4", "gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"},
			Configured:  f.configs["openai"].APIKey != "",
		},
		{
			Name:        "claude",
			DisplayName: "Anthropic Claude",
			Models:      []string{"claude-3-5-sonnet-20241022", "claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307"},
			Configured:  f.configs["claude"].APIKey != "",
		},
		{
			Name:        "groq",
			DisplayName: "Groq",
			Models:      []string{"llama-3.3-70b-versatile", "llama-3.1-70b-versatile", "llama-3.1-8b-instant", "mixtral-8x7b-32768"},
			Configured:  f.configs["groq"].APIKey != "",
		},
		{
			Name:        "gemini",
			DisplayName: "Google Gemini",
			Models:      []string{"gemini-2.0-flash-exp", "gemini-1.5-pro", "gemini-1.5-flash"},
			Configured:  f.configs["gemini"].APIKey != "",
		},
	}
	return infos
}
