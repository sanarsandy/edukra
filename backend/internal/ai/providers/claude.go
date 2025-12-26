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

// ClaudeProvider implements Provider interface for Anthropic Claude
type ClaudeProvider struct {
	config     ProviderConfig
	httpClient *http.Client
	baseURL    string
}

// Claude API types
type claudeChatRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []claudeMessage `json:"messages"`
	System      string          `json:"system,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeChatResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Role         string `json:"role"`
	Content      []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type claudeStreamEvent struct {
	Type         string `json:"type"`
	Index        int    `json:"index,omitempty"`
	ContentBlock *struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content_block,omitempty"`
	Delta *struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	} `json:"delta,omitempty"`
	Message *claudeChatResponse `json:"message,omitempty"`
}

type claudeErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(config ProviderConfig) *ClaudeProvider {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}

	return &ClaudeProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// GetModel returns the current model
func (p *ClaudeProvider) GetModel() string {
	if p.config.Model == "" {
		return "claude-3-5-sonnet-20241022"
	}
	return p.config.Model
}

// AvailableModels returns available Claude models
func (p *ClaudeProvider) AvailableModels() []string {
	return []string{
		"claude-3-5-sonnet-20241022",
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
	}
}

// ValidateAPIKey validates the API key
func (p *ClaudeProvider) ValidateAPIKey() error {
	// Claude doesn't have a simple validation endpoint
	// We'll try a minimal request
	req := ChatRequest{
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
		MaxTokens: 10,
	}

	_, err := p.Chat(context.Background(), req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid_api_key") || strings.Contains(err.Error(), "authentication") {
			return fmt.Errorf("invalid API key")
		}
		return err
	}
	return nil
}

// Chat performs a synchronous chat completion
func (p *ClaudeProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Extract system message and convert messages
	var systemPrompt string
	var messages []claudeMessage

	for _, m := range req.Messages {
		if m.Role == "system" {
			systemPrompt = m.Content
		} else {
			messages = append(messages, claudeMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}
	}

	// Ensure there's at least one message
	if len(messages) == 0 {
		messages = []claudeMessage{{Role: "user", Content: "Hello"}}
	}

	// Prepare request
	claudeReq := claudeChatRequest{
		Model:       p.GetModel(),
		Messages:    messages,
		System:      systemPrompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false,
	}

	// Set defaults
	if claudeReq.MaxTokens == 0 {
		claudeReq.MaxTokens = p.config.MaxTokens
		if claudeReq.MaxTokens == 0 {
			claudeReq.MaxTokens = 2048
		}
	}
	if claudeReq.Temperature == 0 {
		claudeReq.Temperature = p.config.Temperature
		if claudeReq.Temperature == 0 {
			claudeReq.Temperature = 0.7
		}
	}

	body, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Claude: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		var errResp claudeErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("Claude error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("Claude error: %s", string(respBody))
	}

	var claudeResp claudeChatResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return nil, fmt.Errorf("no response from Claude")
	}

	// Combine all text content
	var content strings.Builder
	for _, c := range claudeResp.Content {
		if c.Type == "text" {
			content.WriteString(c.Text)
		}
	}

	return &ChatResponse{
		Content:      content.String(),
		TokensInput:  claudeResp.Usage.InputTokens,
		TokensOutput: claudeResp.Usage.OutputTokens,
		Model:        claudeResp.Model,
		FinishReason: claudeResp.StopReason,
	}, nil
}

// ChatStream performs a streaming chat completion
func (p *ClaudeProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	// Extract system message and convert messages
	var systemPrompt string
	var messages []claudeMessage

	for _, m := range req.Messages {
		if m.Role == "system" {
			systemPrompt = m.Content
		} else {
			messages = append(messages, claudeMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}
	}

	if len(messages) == 0 {
		messages = []claudeMessage{{Role: "user", Content: "Hello"}}
	}

	// Prepare request
	claudeReq := claudeChatRequest{
		Model:       p.GetModel(),
		Messages:    messages,
		System:      systemPrompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      true,
	}

	// Set defaults
	if claudeReq.MaxTokens == 0 {
		claudeReq.MaxTokens = p.config.MaxTokens
		if claudeReq.MaxTokens == 0 {
			claudeReq.MaxTokens = 2048
		}
	}

	body, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Claude: %w", err)
	}

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var errResp claudeErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("Claude error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("Claude error: %s", string(respBody))
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
				if line == "" {
					continue
				}

				if strings.HasPrefix(line, "data: ") {
					jsonData := strings.TrimPrefix(line, "data: ")
					var event claudeStreamEvent
					if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
						continue
					}

					switch event.Type {
					case "content_block_delta":
						if event.Delta != nil && event.Delta.Text != "" {
							ch <- StreamChunk{Content: event.Delta.Text}
						}
					case "message_stop":
						ch <- StreamChunk{Done: true}
						return
					}
				}
			}
		}
	}()

	return ch, nil
}
