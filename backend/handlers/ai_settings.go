package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/providers"
)

// AISettings represents AI configuration
type AISettings struct {
	Enabled             bool    `json:"enabled"`
	Provider            string  `json:"provider"`
	Model               string  `json:"model"`
	APIKeyOpenAI        string  `json:"api_key_openai,omitempty"`
	APIKeyClaude        string  `json:"api_key_claude,omitempty"`
	APIKeyGroq          string  `json:"api_key_groq,omitempty"`
	APIKeyGemini        string  `json:"api_key_gemini,omitempty"`
	EmbeddingProvider   string  `json:"embedding_provider"`
	EmbeddingModel      string  `json:"embedding_model"`
	MaxTokens           int     `json:"max_tokens"`
	Temperature         float64 `json:"temperature"`
	RateLimitPerDay     int     `json:"rate_limit_per_day"`
	SystemPrompt        string  `json:"system_prompt"`
	
	// Read-only status fields
	OpenAIConfigured    bool    `json:"openai_configured"`
	ClaudeConfigured    bool    `json:"claude_configured"`
	GroqConfigured      bool    `json:"groq_configured"`
	GeminiConfigured    bool    `json:"gemini_configured"`
}

// AIUpdateSettingsRequest for updating AI settings
type AIUpdateSettingsRequest struct {
	Enabled             *bool    `json:"enabled,omitempty"`
	Provider            *string  `json:"provider,omitempty"`
	Model               *string  `json:"model,omitempty"`
	APIKeyOpenAI        *string  `json:"api_key_openai,omitempty"`
	APIKeyClaude        *string  `json:"api_key_claude,omitempty"`
	APIKeyGroq          *string  `json:"api_key_groq,omitempty"`
	APIKeyGemini        *string  `json:"api_key_gemini,omitempty"`
	EmbeddingProvider   *string  `json:"embedding_provider,omitempty"`
	EmbeddingModel      *string  `json:"embedding_model,omitempty"`
	MaxTokens           *int     `json:"max_tokens,omitempty"`
	Temperature         *float64 `json:"temperature,omitempty"`
	RateLimitPerDay     *int     `json:"rate_limit_per_day,omitempty"`
	SystemPrompt        *string  `json:"system_prompt,omitempty"`
}

// Simple encryption key - MUST be exactly 32 bytes for AES-256
var encryptionKey = []byte("ai-settings-key-32-bytes-long!!!")

// GetAISettings returns AI configuration (admin only)
func GetAISettings(c echo.Context) error {
	settings := AISettings{
		Enabled:           getSettingBool("ai_enabled", false),
		Provider:          getSettingValue("ai_provider", "openai"),
		Model:             getSettingValue("ai_model", "gpt-4-turbo"),
		EmbeddingProvider: getSettingValue("ai_embedding_provider", "openai"),
		EmbeddingModel:    getSettingValue("ai_embedding_model", "text-embedding-ada-002"),
		MaxTokens:         getSettingInt("ai_max_tokens", 2048),
		Temperature:       getSettingFloat("ai_temperature", 0.7),
		RateLimitPerDay:   getSettingInt("ai_rate_limit_per_day", 50),
		SystemPrompt:      getSettingValue("ai_system_prompt", defaultAISystemPrompt),
	}

	// Check which providers are configured (have API keys)
	settings.OpenAIConfigured = getSettingValue("ai_api_key_openai", "") != ""
	settings.ClaudeConfigured = getSettingValue("ai_api_key_claude", "") != ""
	settings.GroqConfigured = getSettingValue("ai_api_key_groq", "") != ""
	settings.GeminiConfigured = getSettingValue("ai_api_key_gemini", "") != ""

	// Don't return actual API keys, just masked versions
	if settings.OpenAIConfigured {
		settings.APIKeyOpenAI = "sk-****" + maskKey(getSettingValue("ai_api_key_openai", ""))
	}
	if settings.ClaudeConfigured {
		settings.APIKeyClaude = "****" + maskKey(getSettingValue("ai_api_key_claude", ""))
	}
	if settings.GroqConfigured {
		settings.APIKeyGroq = "gsk_****" + maskKey(getSettingValue("ai_api_key_groq", ""))
	}
	if settings.GeminiConfigured {
		settings.APIKeyGemini = "****" + maskKey(getSettingValue("ai_api_key_gemini", ""))
	}

	return c.JSON(http.StatusOK, settings)
}

// UpdateAISettings updates AI configuration (admin only)
func UpdateAISettings(c echo.Context) error {
	var req AIUpdateSettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Debug: Log what we received
	if req.APIKeyGroq != nil {
		log.Printf("[AI Settings] Received Groq API key: length=%d, value=%s...", len(*req.APIKeyGroq), (*req.APIKeyGroq)[:min(10, len(*req.APIKeyGroq))])
	}
	if req.APIKeyOpenAI != nil {
		log.Printf("[AI Settings] Received OpenAI API key: length=%d", len(*req.APIKeyOpenAI))
	}

	// Update each setting if provided
	if req.Enabled != nil {
		setSettingValue("ai_enabled", boolToString(*req.Enabled))
	}
	if req.Provider != nil {
		setSettingValue("ai_provider", *req.Provider)
	}
	if req.Model != nil {
		setSettingValue("ai_model", *req.Model)
	}
	if req.EmbeddingProvider != nil {
		setSettingValue("ai_embedding_provider", *req.EmbeddingProvider)
	}
	if req.EmbeddingModel != nil {
		setSettingValue("ai_embedding_model", *req.EmbeddingModel)
	}
	if req.MaxTokens != nil {
		setSettingValue("ai_max_tokens", strconv.Itoa(*req.MaxTokens))
	}
	if req.Temperature != nil {
		setSettingValue("ai_temperature", strconv.FormatFloat(*req.Temperature, 'f', 2, 64))
	}
	if req.RateLimitPerDay != nil {
		setSettingValue("ai_rate_limit_per_day", strconv.Itoa(*req.RateLimitPerDay))
	}
	if req.SystemPrompt != nil {
		setSettingValue("ai_system_prompt", *req.SystemPrompt)
	}

	// Update API keys (encrypt before storing)
	if req.APIKeyOpenAI != nil && *req.APIKeyOpenAI != "" && !isPlaceholder(*req.APIKeyOpenAI) {
		encrypted, err := encrypt(*req.APIKeyOpenAI)
		if err == nil {
			setSettingValue("ai_api_key_openai", encrypted)
		}
	}
	if req.APIKeyClaude != nil && *req.APIKeyClaude != "" && !isPlaceholder(*req.APIKeyClaude) {
		encrypted, err := encrypt(*req.APIKeyClaude)
		if err == nil {
			setSettingValue("ai_api_key_claude", encrypted)
		}
	}
	if req.APIKeyGroq != nil && *req.APIKeyGroq != "" && !isPlaceholder(*req.APIKeyGroq) {
		log.Printf("[AI Settings] Groq key passed checks, encrypting...")
		encrypted, err := encrypt(*req.APIKeyGroq)
		if err != nil {
			log.Printf("[AI Settings] Encryption FAILED: %v", err)
		} else {
			log.Printf("[AI Settings] Encrypted key length: %d, saving to DB...", len(encrypted))
			saveErr := setSettingValue("ai_api_key_groq", encrypted)
			if saveErr != nil {
				log.Printf("[AI Settings] Save FAILED: %v", saveErr)
			} else {
				log.Printf("[AI Settings] Save SUCCESS!")
			}
		}
	} else if req.APIKeyGroq != nil {
		log.Printf("[AI Settings] Groq key rejected: empty=%v, isPlaceholder=%v", *req.APIKeyGroq == "", isPlaceholder(*req.APIKeyGroq))
	}
	if req.APIKeyGemini != nil && *req.APIKeyGemini != "" && !isPlaceholder(*req.APIKeyGemini) {
		encrypted, err := encrypt(*req.APIKeyGemini)
		if err == nil {
			setSettingValue("ai_api_key_gemini", encrypted)
		}
	}

	return GetAISettings(c)
}

// ValidateAIKey validates an API key for a provider
func ValidateAIKey(c echo.Context) error {
	var req struct {
		Provider string `json:"provider"`
		APIKey   string `json:"api_key"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Provider == "" || req.APIKey == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Provider and API key required"})
	}

	config := providers.ProviderConfig{
		APIKey: req.APIKey,
	}

	var provider providers.Provider
	switch req.Provider {
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

	if err := provider.ValidateAPIKey(); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"valid":   false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"valid":   true,
		"message": "API key is valid",
	})
}

// GetAIProviders returns available AI providers and their models
func GetAIProviders(c echo.Context) error {
	factory := providers.NewProviderFactory()
	infos := factory.GetProviderInfo()
	return c.JSON(http.StatusOK, infos)
}

// ClearAPIKey removes an API key for a provider
func ClearAPIKey(c echo.Context) error {
	provider := c.QueryParam("provider")
	
	var key string
	switch provider {
	case "openai":
		key = "ai_api_key_openai"
	case "claude":
		key = "ai_api_key_claude"
	case "groq":
		key = "ai_api_key_groq"
	case "gemini":
		key = "ai_api_key_gemini"
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unknown provider"})
	}

	setSettingValue(key, "")
	return c.JSON(http.StatusOK, map[string]string{"message": "API key cleared"})
}

// Helper functions

func getSettingBool(key string, defaultVal bool) bool {
	val := getSettingValue(key, "")
	if val == "" {
		return defaultVal
	}
	return val == "true" || val == "1"
}

func getSettingInt(key string, defaultVal int) int {
	val := getSettingValue(key, "")
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func getSettingFloat(key string, defaultVal float64) float64 {
	val := getSettingValue(key, "")
	if val == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func maskKey(key string) string {
	if len(key) <= 4 {
		return "****"
	}
	return key[len(key)-4:]
}

func isPlaceholder(key string) bool {
	// A placeholder is a masked key shown in the UI like "****xxxx" or "sk-****xxxx"
	if len(key) < 8 {
		return true // Too short to be a real key
	}
	// Check if key contains ****
	if len(key) >= 4 && key[0:4] == "****" {
		return true
	}
	// Check for patterns like "sk-****" or "gsk_****"
	for i := 0; i < len(key)-4; i++ {
		if key[i:i+4] == "****" {
			return true
		}
	}
	return false
}

func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GetDecryptedAPIKey returns decrypted API key for a provider
func GetDecryptedAPIKey(provider string) (string, error) {
	var key string
	switch provider {
	case "openai":
		key = getSettingValue("ai_api_key_openai", "")
	case "claude":
		key = getSettingValue("ai_api_key_claude", "")
	case "groq":
		key = getSettingValue("ai_api_key_groq", "")
	case "gemini":
		key = getSettingValue("ai_api_key_gemini", "")
	default:
		return "", nil
	}

	if key == "" {
		return "", nil
	}

	return decrypt(key)
}

const defaultAISystemPrompt = `Kamu adalah AI Tutor yang membantu siswa memahami materi kursus. 
Tugasmu adalah:
1. Menjawab pertanyaan siswa berdasarkan materi yang diberikan
2. Memberikan penjelasan yang jelas dan mudah dipahami
3. Menggunakan bahasa Indonesia yang baik dan benar
4. Mengutip sumber materi ketika relevan

Aturan penting:
- Jika pertanyaan tidak bisa dijawab dari materi yang tersedia, katakan dengan jujur bahwa kamu tidak memiliki informasi tersebut
- Jangan mengarang jawaban yang tidak ada dalam materi
- Berikan jawaban yang ringkas namun lengkap
- Gunakan format yang mudah dibaca (bullet points, numbering) jika sesuai`

// Ensure unused imports are used
var _ = db.DB
var _ = json.Marshal

// FetchProviderModels fetches available models from a provider's API
func FetchProviderModels(c echo.Context) error {
	provider := c.QueryParam("provider")
	if provider == "" {
		provider = "gemini"
	}

	apiKey, err := GetDecryptedAPIKey(provider)
	if err != nil || apiKey == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "API key tidak ditemukan untuk provider " + provider,
			"models": []string{},
		})
	}

	var models []map[string]interface{}
	var rawResponse string

	switch provider {
	case "gemini":
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s", apiKey)
		resp, err := http.Get(url)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":  "Gagal menghubungi Gemini API: " + err.Error(),
				"models": []string{},
			})
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		rawResponse = string(body)

		if resp.StatusCode != 200 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":        "Gemini API error",
				"raw_response": rawResponse,
				"models":       []string{},
			})
		}

		var result struct {
			Models []struct {
				Name                       string   `json:"name"`
				DisplayName                string   `json:"displayName"`
				Description                string   `json:"description"`
				SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
			} `json:"models"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":        "Gagal parse response: " + err.Error(),
				"raw_response": rawResponse,
				"models":       []string{},
			})
		}

		// Filter models that support generateContent
		for _, m := range result.Models {
			supportsGenerateContent := false
			supportsEmbedContent := false
			for _, method := range m.SupportedGenerationMethods {
				if method == "generateContent" {
					supportsGenerateContent = true
				}
				if method == "embedContent" {
					supportsEmbedContent = true
				}
			}
			// Extract model name from full path (models/gemini-1.5-pro -> gemini-1.5-pro)
			modelName := m.Name
			if len(m.Name) > 7 && m.Name[:7] == "models/" {
				modelName = m.Name[7:]
			}
			models = append(models, map[string]interface{}{
				"name":                   modelName,
				"displayName":            m.DisplayName,
				"description":            m.Description,
				"supportsGenerateContent": supportsGenerateContent,
				"supportsEmbedContent":   supportsEmbedContent,
			})
		}
	default:
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Provider tidak didukung untuk fetch models: " + provider,
			"models": []string{},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"provider": provider,
		"models":   models,
		"count":    len(models),
	})
}
