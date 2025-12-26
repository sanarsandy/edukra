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

// EmbeddingVector represents a vector embedding
type EmbeddingVector []float32

// EmbeddingResult contains the embedding and metadata
type EmbeddingResult struct {
	Embedding   EmbeddingVector `json:"embedding"`
	TokensUsed  int             `json:"tokens_used"`
	Model       string          `json:"model"`
}

// Embedder interface for generating embeddings
type Embedder interface {
	// Name returns the embedder name
	Name() string

	// Embed generates embedding for a single text
	Embed(ctx context.Context, text string) (*EmbeddingResult, error)

	// EmbedBatch generates embeddings for multiple texts
	EmbedBatch(ctx context.Context, texts []string) ([]EmbeddingResult, error)

	// Dimensions returns the embedding dimensions
	Dimensions() int

	// GetModel returns the current model
	GetModel() string
}

// EmbedderConfig holds configuration for embedders
type EmbedderConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

// OpenAIEmbedder implements Embedder using OpenAI
type OpenAIEmbedder struct {
	config     EmbedderConfig
	httpClient *http.Client
	baseURL    string
}

// OpenAI embedding API types
type openAIEmbeddingRequest struct {
	Input          interface{} `json:"input"` // string or []string
	Model          string      `json:"model"`
	EncodingFormat string      `json:"encoding_format,omitempty"`
}

type openAIEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type openAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// NewOpenAIEmbedder creates a new OpenAI embedder
func NewOpenAIEmbedder(config EmbedderConfig) *OpenAIEmbedder {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIEmbedder{
		config: config,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Name returns the embedder name
func (e *OpenAIEmbedder) Name() string {
	return "openai"
}

// GetModel returns the current model
func (e *OpenAIEmbedder) GetModel() string {
	if e.config.Model == "" {
		return "text-embedding-ada-002"
	}
	return e.config.Model
}

// Dimensions returns embedding dimensions based on model
func (e *OpenAIEmbedder) Dimensions() int {
	model := e.GetModel()
	switch model {
	case "text-embedding-3-large":
		return 3072
	case "text-embedding-3-small":
		return 1536
	case "text-embedding-ada-002":
		return 1536
	default:
		return 1536
	}
}

// Embed generates embedding for a single text
func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) (*EmbeddingResult, error) {
	results, err := e.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return &results[0], nil
}

// EmbedBatch generates embeddings for multiple texts
func (e *OpenAIEmbedder) EmbedBatch(ctx context.Context, texts []string) ([]EmbeddingResult, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	// Prepare request
	reqBody := openAIEmbeddingRequest{
		Input: texts,
		Model: e.GetModel(),
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", e.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+e.config.APIKey)

	resp, err := e.httpClient.Do(httpReq)
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

	var embResp openAIEmbeddingResponse
	if err := json.Unmarshal(respBody, &embResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Sort results by index and build response
	results := make([]EmbeddingResult, len(embResp.Data))
	tokensPerItem := embResp.Usage.TotalTokens / len(texts)
	
	for _, data := range embResp.Data {
		results[data.Index] = EmbeddingResult{
			Embedding:  data.Embedding,
			TokensUsed: tokensPerItem,
			Model:      embResp.Model,
		}
	}

	return results, nil
}

// AvailableEmbeddingModels returns list of available embedding models
func AvailableEmbeddingModels() []string {
	return []string{
		"text-embedding-ada-002",
		"text-embedding-3-small",
		"text-embedding-3-large",
	}
}
