package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GeminiEmbedder implements Embedder using Google Gemini
type GeminiEmbedder struct {
	config     EmbedderConfig
	httpClient *http.Client
	baseURL    string
}

// Gemini embedding API types
type geminiEmbedRequest struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type geminiEmbedResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
}

type geminiBatchEmbedRequest struct {
	Requests []geminiEmbedRequest `json:"requests"`
}

type geminiBatchEmbedResponse struct {
	Embeddings []struct {
		Values []float32 `json:"values"`
	} `json:"embeddings"`
}

type geminiErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// NewGeminiEmbedder creates a new Gemini embedder
func NewGeminiEmbedder(config EmbedderConfig) *GeminiEmbedder {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	if config.Model == "" {
		config.Model = "gemini-embedding-001"
	}

	return &GeminiEmbedder{
		config: config,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Name returns the embedder name
func (e *GeminiEmbedder) Name() string {
	return "gemini"
}

// GetModel returns the current model
func (e *GeminiEmbedder) GetModel() string {
	if e.config.Model == "" {
		return "gemini-embedding-001"
	}
	return e.config.Model
}

// Dimensions returns embedding dimensions based on model
func (e *GeminiEmbedder) Dimensions() int {
	// Gemini text-embedding-004 returns 768-dimensional embeddings
	return 768
}

// Embed generates embedding for a single text
func (e *GeminiEmbedder) Embed(ctx context.Context, text string) (*EmbeddingResult, error) {
	model := e.GetModel()
	url := fmt.Sprintf("%s/models/%s:embedContent?key=%s", e.baseURL, model, e.config.APIKey)

	reqBody := geminiEmbedRequest{}
	reqBody.Content.Parts = []struct {
		Text string `json:"text"`
	}{{Text: text}}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(httpReq)
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

	var embResp geminiEmbedResponse
	if err := json.Unmarshal(respBody, &embResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &EmbeddingResult{
		Embedding:  embResp.Embedding.Values,
		TokensUsed: len(text) / 4, // Rough estimate
		Model:      model,
	}, nil
}

// EmbedBatch generates embeddings for multiple texts
func (e *GeminiEmbedder) EmbedBatch(ctx context.Context, texts []string) ([]EmbeddingResult, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	// Gemini doesn't have a batch endpoint as straightforward as OpenAI
	// So we process one at a time (could be optimized with concurrent calls)
	results := make([]EmbeddingResult, len(texts))
	for i, text := range texts {
		result, err := e.Embed(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to embed text %d: %w", i, err)
		}
		results[i] = *result
	}

	return results, nil
}

// AvailableGeminiEmbeddingModels returns list of available Gemini embedding models
func AvailableGeminiEmbeddingModels() []string {
	return []string{
		"gemini-embedding-001",
		"text-embedding-004",
	}
}
