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

// GroqProvider implements Provider interface for Groq
type GroqProvider struct {
	config     ProviderConfig
	httpClient *http.Client
	baseURL    string
}

// Groq uses OpenAI-compatible API format
type groqChatRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      groqMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type groqStreamChunk struct {
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

type groqErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// NewGroqProvider creates a new Groq provider
func NewGroqProvider(config ProviderConfig) *GroqProvider {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.groq.com/openai/v1"
	}

	return &GroqProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Name returns the provider name
func (p *GroqProvider) Name() string {
	return "groq"
}

// GetModel returns the current model
func (p *GroqProvider) GetModel() string {
	if p.config.Model == "" {
		return "llama-3.3-70b-versatile"
	}
	return p.config.Model
}

// AvailableModels returns available Groq models
func (p *GroqProvider) AvailableModels() []string {
	return []string{
		"llama-3.3-70b-versatile",
		"llama-3.1-70b-versatile",
		"llama-3.1-8b-instant",
		"mixtral-8x7b-32768",
		"gemma2-9b-it",
	}
}

// ValidateAPIKey validates the API key
func (p *GroqProvider) ValidateAPIKey() error {
	req, err := http.NewRequest("GET", p.baseURL+"/models", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Groq: %w", err)
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
func (p *GroqProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Convert messages
	messages := make([]groqMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = groqMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	// Prepare request
	groqReq := groqChatRequest{
		Model:       p.GetModel(),
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false,
	}

	// Set defaults
	if groqReq.MaxTokens == 0 {
		groqReq.MaxTokens = p.config.MaxTokens
		if groqReq.MaxTokens == 0 {
			groqReq.MaxTokens = 2048
		}
	}
	if groqReq.Temperature == 0 {
		groqReq.Temperature = p.config.Temperature
		if groqReq.Temperature == 0 {
			groqReq.Temperature = 0.7
		}
	}

	body, err := json.Marshal(groqReq)
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
		return nil, fmt.Errorf("failed to call Groq: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		var errResp groqErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("Groq error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("Groq error: %s", string(respBody))
	}

	var groqResp groqChatResponse
	if err := json.Unmarshal(respBody, &groqResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from Groq")
	}

	return &ChatResponse{
		Content:      groqResp.Choices[0].Message.Content,
		TokensInput:  groqResp.Usage.PromptTokens,
		TokensOutput: groqResp.Usage.CompletionTokens,
		Model:        groqResp.Model,
		FinishReason: groqResp.Choices[0].FinishReason,
	}, nil
}

// ChatStream performs a streaming chat completion
func (p *GroqProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	// Convert messages
	messages := make([]groqMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = groqMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	// Prepare request
	groqReq := groqChatRequest{
		Model:       p.GetModel(),
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      true,
	}

	// Set defaults
	if groqReq.MaxTokens == 0 {
		groqReq.MaxTokens = p.config.MaxTokens
		if groqReq.MaxTokens == 0 {
			groqReq.MaxTokens = 2048
		}
	}

	body, err := json.Marshal(groqReq)
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
		return nil, fmt.Errorf("failed to call Groq: %w", err)
	}

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var errResp groqErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("Groq error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("Groq error: %s", string(respBody))
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
					var chunk groqStreamChunk
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
