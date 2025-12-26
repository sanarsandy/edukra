package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAIProvider implements Provider interface for OpenAI
type OpenAIProvider struct {
	config     ProviderConfig
	httpClient *http.Client
	baseURL    string
}

// OpenAI API types
type openAIChatRequest struct {
	Model       string           `json:"model"`
	Messages    []openAIMessage  `json:"messages"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
	Temperature float64          `json:"temperature,omitempty"`
	Stream      bool             `json:"stream,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      openAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type openAIStreamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

type openAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config ProviderConfig) *OpenAIProvider {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// GetModel returns the current model
func (p *OpenAIProvider) GetModel() string {
	if p.config.Model == "" {
		return "gpt-4-turbo"
	}
	return p.config.Model
}

// AvailableModels returns available OpenAI models
func (p *OpenAIProvider) AvailableModels() []string {
	return []string{
		"gpt-4-turbo",
		"gpt-4",
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-3.5-turbo",
	}
}

// ValidateAPIKey validates the API key
func (p *OpenAIProvider) ValidateAPIKey() error {
	req, err := http.NewRequest("GET", p.baseURL+"/models", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}

	return nil
}

// Chat performs a synchronous chat completion
func (p *OpenAIProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Convert messages
	messages := make([]openAIMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openAIMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	// Prepare request
	openAIReq := openAIChatRequest{
		Model:       p.GetModel(),
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false,
	}

	// Set defaults
	if openAIReq.MaxTokens == 0 {
		openAIReq.MaxTokens = p.config.MaxTokens
		if openAIReq.MaxTokens == 0 {
			openAIReq.MaxTokens = 2048
		}
	}
	if openAIReq.Temperature == 0 {
		openAIReq.Temperature = p.config.Temperature
		if openAIReq.Temperature == 0 {
			openAIReq.Temperature = 0.7
		}
	}

	body, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		var errResp openAIErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("OpenAI error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("OpenAI error: %s", string(respBody))
	}

	var openAIResp openAIChatResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return &ChatResponse{
		Content:      openAIResp.Choices[0].Message.Content,
		TokensInput:  openAIResp.Usage.PromptTokens,
		TokensOutput: openAIResp.Usage.CompletionTokens,
		Model:        openAIResp.Model,
		FinishReason: openAIResp.Choices[0].FinishReason,
	}, nil
}

// ChatStream performs a streaming chat completion
func (p *OpenAIProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	// Convert messages
	messages := make([]openAIMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openAIMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	// Prepare request
	openAIReq := openAIChatRequest{
		Model:       p.GetModel(),
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      true,
	}

	// Set defaults
	if openAIReq.MaxTokens == 0 {
		openAIReq.MaxTokens = p.config.MaxTokens
		if openAIReq.MaxTokens == 0 {
			openAIReq.MaxTokens = 2048
		}
	}
	if openAIReq.Temperature == 0 {
		openAIReq.Temperature = p.config.Temperature
		if openAIReq.Temperature == 0 {
			openAIReq.Temperature = 0.7
		}
	}

	body, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI: %w", err)
	}

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var errResp openAIErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("OpenAI error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("OpenAI error: %s", string(respBody))
	}

	ch := make(chan StreamChunk, 100)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		reader := resp.Body
		buffer := make([]byte, 4096)

		for {
			select {
			case <-ctx.Done():
				ch <- StreamChunk{Done: true, Error: ctx.Err()}
				return
			default:
			}

			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					ch <- StreamChunk{Done: true}
					return
				}
				ch <- StreamChunk{Done: true, Error: err}
				return
			}

			data := string(buffer[:n])
			lines := strings.Split(data, "\n")

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || line == "data: [DONE]" {
					continue
				}

				if strings.HasPrefix(line, "data: ") {
					jsonData := strings.TrimPrefix(line, "data: ")
					var chunk openAIStreamChunk
					if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
						continue
					}

					if len(chunk.Choices) > 0 {
						content := chunk.Choices[0].Delta.Content
						if content != "" {
							ch <- StreamChunk{Content: content}
						}

						if chunk.Choices[0].FinishReason != nil {
							ch <- StreamChunk{Done: true}
							return
						}
					}
				}
			}
		}
	}()

	return ch, nil
}
