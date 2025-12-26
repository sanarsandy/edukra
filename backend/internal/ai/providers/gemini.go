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

// GeminiProvider implements Provider interface for Google Gemini
type GeminiProvider struct {
	config     ProviderConfig
	httpClient *http.Client
	baseURL    string
}

// Gemini API types
type geminiGenerateRequest struct {
	Contents         []geminiContent       `json:"contents"`
	SystemInstruction *geminiContent       `json:"systemInstruction,omitempty"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
	Temperature     float64 `json:"temperature,omitempty"`
}

type geminiGenerateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

type geminiErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(config ProviderConfig) *GeminiProvider {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	return &GeminiProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Name returns the provider name
func (p *GeminiProvider) Name() string {
	return "gemini"
}

// GetModel returns the current model
func (p *GeminiProvider) GetModel() string {
	if p.config.Model == "" {
		return "gemini-2.0-flash-exp"
	}
	return p.config.Model
}

// AvailableModels returns available Gemini models
func (p *GeminiProvider) AvailableModels() []string {
	return []string{
		"gemini-3-flash",
		"gemini-2.5-flash",
		"gemini-2.5-flash-lite",
	}
}

// ValidateAPIKey validates the API key
func (p *GeminiProvider) ValidateAPIKey() error {
	url := fmt.Sprintf("%s/models?key=%s", p.baseURL, p.config.APIKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Gemini: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 || resp.StatusCode == 401 || resp.StatusCode == 403 {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}

	return nil
}

// Chat performs a synchronous chat completion
func (p *GeminiProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Extract system message and convert messages
	var systemInstruction *geminiContent
	var contents []geminiContent

	for _, m := range req.Messages {
		if m.Role == "system" {
			systemInstruction = &geminiContent{
				Parts: []geminiPart{{Text: m.Content}},
			}
		} else {
			role := m.Role
			if role == "assistant" {
				role = "model"
			}
			contents = append(contents, geminiContent{
				Parts: []geminiPart{{Text: m.Content}},
				Role:  role,
			})
		}
	}

	// Ensure there's at least one message
	if len(contents) == 0 {
		contents = []geminiContent{{
			Parts: []geminiPart{{Text: "Hello"}},
			Role:  "user",
		}}
	}

	// Prepare request
	geminiReq := geminiGenerateRequest{
		Contents:         contents,
		SystemInstruction: systemInstruction,
		GenerationConfig: geminiGenerationConfig{
			MaxOutputTokens: req.MaxTokens,
			Temperature:     req.Temperature,
		},
	}

	// Set defaults
	if geminiReq.GenerationConfig.MaxOutputTokens == 0 {
		geminiReq.GenerationConfig.MaxOutputTokens = p.config.MaxTokens
		if geminiReq.GenerationConfig.MaxOutputTokens == 0 {
			geminiReq.GenerationConfig.MaxOutputTokens = 2048
		}
	}
	if geminiReq.GenerationConfig.Temperature == 0 {
		geminiReq.GenerationConfig.Temperature = p.config.Temperature
		if geminiReq.GenerationConfig.Temperature == 0 {
			geminiReq.GenerationConfig.Temperature = 0.7
		}
	}

	body, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.baseURL, p.GetModel(), p.config.APIKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		var errResp geminiErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("Gemini error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("Gemini error: %s", string(respBody))
	}

	var geminiResp geminiGenerateResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	// Combine all text parts
	var content strings.Builder
	for _, part := range geminiResp.Candidates[0].Content.Parts {
		content.WriteString(part.Text)
	}

	return &ChatResponse{
		Content:      content.String(),
		TokensInput:  geminiResp.UsageMetadata.PromptTokenCount,
		TokensOutput: geminiResp.UsageMetadata.CandidatesTokenCount,
		Model:        p.GetModel(),
		FinishReason: geminiResp.Candidates[0].FinishReason,
	}, nil
}

// ChatStream performs a streaming chat completion
func (p *GeminiProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	// Extract system message and convert messages
	var systemInstruction *geminiContent
	var contents []geminiContent

	for _, m := range req.Messages {
		if m.Role == "system" {
			systemInstruction = &geminiContent{
				Parts: []geminiPart{{Text: m.Content}},
			}
		} else {
			role := m.Role
			if role == "assistant" {
				role = "model"
			}
			contents = append(contents, geminiContent{
				Parts: []geminiPart{{Text: m.Content}},
				Role:  role,
			})
		}
	}

	if len(contents) == 0 {
		contents = []geminiContent{{
			Parts: []geminiPart{{Text: "Hello"}},
			Role:  "user",
		}}
	}

	// Prepare request
	geminiReq := geminiGenerateRequest{
		Contents:         contents,
		SystemInstruction: systemInstruction,
		GenerationConfig: geminiGenerationConfig{
			MaxOutputTokens: req.MaxTokens,
			Temperature:     req.Temperature,
		},
	}

	// Set defaults
	if geminiReq.GenerationConfig.MaxOutputTokens == 0 {
		geminiReq.GenerationConfig.MaxOutputTokens = p.config.MaxTokens
		if geminiReq.GenerationConfig.MaxOutputTokens == 0 {
			geminiReq.GenerationConfig.MaxOutputTokens = 2048
		}
	}

	body, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:streamGenerateContent?alt=sse&key=%s", p.baseURL, p.GetModel(), p.config.APIKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini: %w", err)
	}

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var errResp geminiErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("Gemini error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("Gemini error: %s", string(respBody))
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
					var geminiResp geminiGenerateResponse
					if err := json.Unmarshal([]byte(jsonData), &geminiResp); err != nil {
						continue
					}

					if len(geminiResp.Candidates) > 0 {
						for _, part := range geminiResp.Candidates[0].Content.Parts {
							if part.Text != "" {
								ch <- StreamChunk{Content: part.Text}
							}
						}

						if geminiResp.Candidates[0].FinishReason != "" {
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
