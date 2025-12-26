package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/embeddings"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/providers"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/rag"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

// Global AI repository (initialized lazily)
var aiRepo *postgres.AIRepository

func getAIRepo() *postgres.AIRepository {
	if aiRepo == nil {
		aiRepo = postgres.NewAIRepository(db.DB)
	}
	return aiRepo
}

// ChatRequest represents a chat message request
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	Message    string       `json:"message"`
	Sources    []rag.Source `json:"sources,omitempty"`
	SessionID  string       `json:"session_id"`
	TokensUsed int          `json:"tokens_used"`
	Quota      QuotaInfo    `json:"quota"`
}

// QuotaInfo shows remaining quota
type QuotaInfo struct {
	Used      int `json:"used"`
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
}

// SendChatMessage handles chat messages from students
func SendChatMessage(c echo.Context) error {
	// Check if AI is enabled
	if !getSettingBool("ai_enabled", false) {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "AI Tutor belum diaktifkan",
		})
	}

	courseID := c.Param("id")
	userID := getUserIDFromToken(c)

	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req ChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Message == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Message is required"})
	}

	ctx := c.Request().Context()
	repo := getAIRepo()

	// Check rate limit
	rateLimit := getSettingInt("ai_rate_limit_per_day", 50)
	todayUsage, _ := repo.GetTodayUsageCount(ctx, userID)
	
	if todayUsage >= rateLimit {
		return c.JSON(http.StatusTooManyRequests, map[string]string{
			"error": "Batas harian tercapai. Silakan coba lagi besok.",
		})
	}

	// Get or create session
	session, err := repo.GetOrCreateActiveSession(ctx, userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get session"})
	}

	// Get AI provider
	providerName := getSettingValue("ai_provider", "openai")
	apiKey, err := GetDecryptedAPIKey(providerName)
	if err != nil || apiKey == "" {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "AI provider belum dikonfigurasi",
		})
	}

	// Get chat history
	history, _ := repo.GetChatHistory(ctx, session.ID, 20)
	
	// Convert to provider messages
	var chatHistory []providers.Message
	for _, msg := range history {
		chatHistory = append(chatHistory, providers.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Save user message
	repo.AddChatMessage(ctx, session.ID, "user", req.Message, nil, 0, "", "")

	// Build prompt with RAG context
	var contextStr string
	var sources []rag.Source

	// Try to get relevant context
	embeddingCount, _ := repo.GetEmbeddingCount(ctx, courseID)
	log.Printf("[AI Chat] Course %s has %d embeddings", courseID, embeddingCount)
	if embeddingCount > 0 {
		log.Printf("[AI Chat] Retrieving RAG context for query: %s", req.Message)
		contextStr, sources = getRAGContext(ctx, courseID, req.Message, apiKey)
		log.Printf("[AI Chat] Got context length: %d, sources: %d", len(contextStr), len(sources))
	} else {
		log.Printf("[AI Chat] No embeddings found, skipping RAG")
	}

	// Build messages
	promptBuilder := rag.NewPromptBuilder(rag.DefaultPromptTemplate())
	
	// Use custom system prompt if set
	customPrompt := getSettingValue("ai_system_prompt", "")
	if customPrompt != "" {
		promptBuilder = rag.NewPromptBuilder(rag.PromptTemplate{
			SystemPrompt:   customPrompt,
			ContextPrompt:  rag.DefaultPromptTemplate().ContextPrompt,
			QuestionPrompt: rag.DefaultPromptTemplate().QuestionPrompt,
		})
	}

	var messages []providers.Message
	if contextStr != "" {
		messages = promptBuilder.BuildMessages(contextStr, req.Message, chatHistory)
	} else {
		messages = promptBuilder.BuildSimpleMessages(req.Message, chatHistory)
	}

	// Call AI provider
	config := providers.ProviderConfig{
		APIKey:      apiKey,
		Model:       getSettingValue("ai_model", "gpt-4-turbo"),
		MaxTokens:   getSettingInt("ai_max_tokens", 2048),
		Temperature: getSettingFloat("ai_temperature", 0.7),
	}

	var provider providers.Provider
	switch providerName {
	case "openai":
		provider = providers.NewOpenAIProvider(config)
	case "claude":
		provider = providers.NewClaudeProvider(config)
	case "groq":
		provider = providers.NewGroqProvider(config)
	case "gemini":
		provider = providers.NewGeminiProvider(config)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unknown provider"})
	}

	response, err := provider.Chat(ctx, providers.ChatRequest{
		Messages:    messages,
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mendapatkan respons AI: " + err.Error(),
		})
	}

	// Save assistant message
	totalTokens := response.TokensInput + response.TokensOutput
	repo.AddChatMessage(ctx, session.ID, "assistant", response.Content, sources, totalTokens, providerName, response.Model)

	// Log usage
	repo.LogUsage(ctx, userID, courseID, "chat", providerName, response.Model, response.TokensInput, response.TokensOutput)

	return c.JSON(http.StatusOK, ChatResponse{
		Message:    response.Content,
		Sources:    sources,
		SessionID:  session.ID,
		TokensUsed: totalTokens,
		Quota: QuotaInfo{
			Used:      todayUsage + 1,
			Limit:     rateLimit,
			Remaining: rateLimit - todayUsage - 1,
		},
	})
}

// GetChatSession returns the current chat session
func GetChatSession(c echo.Context) error {
	courseID := c.Param("id")
	userID := getUserIDFromToken(c)

	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	ctx := c.Request().Context()
	repo := getAIRepo()

	session, err := repo.GetOrCreateActiveSession(ctx, userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get session"})
	}

	messages, _ := repo.GetChatHistory(ctx, session.ID, 100)

	// Check quota
	rateLimit := getSettingInt("ai_rate_limit_per_day", 50)
	todayUsage, _ := repo.GetTodayUsageCount(ctx, userID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"session":  session,
		"messages": messages,
		"quota": QuotaInfo{
			Used:      todayUsage,
			Limit:     rateLimit,
			Remaining: rateLimit - todayUsage,
		},
		"ai_enabled": getSettingBool("ai_enabled", false),
	})
}

// ClearChatSession clears the current chat session
func ClearChatSession(c echo.Context) error {
	courseID := c.Param("id")
	userID := getUserIDFromToken(c)

	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	ctx := c.Request().Context()
	repo := getAIRepo()

	session, err := repo.GetOrCreateActiveSession(ctx, userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get session"})
	}

	if err := repo.ClearChatSession(ctx, session.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to clear session"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Session cleared"})
}

// GetChatQuota returns the user's remaining chat quota
func GetChatQuota(c echo.Context) error {
	userID := getUserIDFromToken(c)

	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	ctx := c.Request().Context()
	repo := getAIRepo()

	rateLimit := getSettingInt("ai_rate_limit_per_day", 50)
	todayUsage, _ := repo.GetTodayUsageCount(ctx, userID)

	return c.JSON(http.StatusOK, QuotaInfo{
		Used:      todayUsage,
		Limit:     rateLimit,
		Remaining: rateLimit - todayUsage,
	})
}

// GetAIStatus returns AI status for a course
func GetAIStatus(c echo.Context) error {
	courseID := c.Param("id")
	ctx := c.Request().Context()
	repo := getAIRepo()

	embeddingCount, _ := repo.GetEmbeddingCount(ctx, courseID)
	processingStatus, _ := repo.GetProcessingStatus(ctx, courseID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ai_enabled":        getSettingBool("ai_enabled", false),
		"embedding_count":   embeddingCount,
		"processing_status": processingStatus,
		"provider":          getSettingValue("ai_provider", "openai"),
		"model":             getSettingValue("ai_model", "gpt-4-turbo"),
	})
}

// Helper to get RAG context
func getRAGContext(ctx context.Context, courseID, query, apiKey string) (string, []rag.Source) {
	repo := getAIRepo()
	
	// Get the AI provider to use correct embedder
	providerName := getSettingValue("ai_provider", "openai")
	embeddingModel := getSettingValue("ai_embedding_model", "")
	
	// Create embedder config
	embConfig := embeddings.EmbedderConfig{
		APIKey: apiKey,
		Model:  embeddingModel,
	}
	
	// Create correct embedder based on provider
	var embedder embeddings.Embedder
	switch providerName {
	case "gemini":
		if embConfig.Model == "" {
			embConfig.Model = "gemini-embedding-001"

		}
		embedder = embeddings.NewGeminiEmbedder(embConfig)
	default: // openai
		if embConfig.Model == "" {
			embConfig.Model = "text-embedding-ada-002"
		}
		embedder = embeddings.NewOpenAIEmbedder(embConfig)
	}

	// Create retriever
	retrieverConfig := rag.DefaultRetrieverConfig()
	retriever := rag.NewRetriever(embedder, repo, retrieverConfig)

	// Retrieve relevant chunks
	chunks, err := retriever.Retrieve(ctx, courseID, query)
	if err != nil || len(chunks) == 0 {
		return "", nil
	}

	// Build context
	return retriever.BuildContext(chunks)
}

