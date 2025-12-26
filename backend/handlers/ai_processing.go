package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/embeddings"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/processor"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/rag"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

// ProcessCourseContent triggers AI content processing for a course
func ProcessCourseContent(c echo.Context) error {
	courseID := c.Param("id")

	// Check if AI is enabled
	if !getSettingBool("ai_enabled", false) {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "AI Tutor belum diaktifkan",
		})
	}

	// Get the selected AI provider
	aiProvider := getSettingValue("ai_provider", "openai")
	
	// Check if provider supports embeddings
	// Currently OpenAI and Gemini have embedding APIs
	supportedProviders := map[string]bool{"openai": true, "gemini": true}
	if !supportedProviders[aiProvider] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Provider " + aiProvider + " tidak mendukung embedding. Pilih OpenAI atau Gemini di Settings AI Tutor untuk menggunakan fitur RAG.",
		})
	}

	// Get API key from the selected AI provider
	apiKey, _ := GetDecryptedAPIKey(aiProvider)
	if apiKey == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "API key belum dikonfigurasi untuk provider " + aiProvider,
		})
	}

	// Create processing status BEFORE starting goroutine (to avoid race condition)
	// This ensures frontend sees "processing" status immediately when polling
	repo := postgres.NewAIRepository(db.DB)
	status, err := repo.CreateProcessingStatus(context.Background(), courseID, "", "all")
	if err != nil {
		log.Printf("[AI Processing] Failed to create status: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal membuat status pemrosesan",
		})
	}
	
	// Update status to "processing" immediately
	if err := repo.UpdateProcessingStatus(context.Background(), status.ID, "processing", 0, 0, nil); err != nil {
		log.Printf("[AI Processing] ERROR updating status to processing: %v", err)
	} else {
		log.Printf("[AI Processing] Created status ID: %s with status=processing", status.ID)
	}

	// Start processing in background with NEW context (not request context)
	go processCourseContentAsync(context.Background(), courseID, apiKey, aiProvider, status.ID)

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "Pemrosesan konten dimulai",
		"status":  "processing",
	})
}

// processCourseContentAsync processes course content in background
// statusID is passed from caller (created before goroutine to avoid race condition)
func processCourseContentAsync(ctx context.Context, courseID, apiKey, provider, statusID string) {
	log.Printf("[AI Processing] Starting for course %s with provider %s, statusID: %s", courseID, provider, statusID)
	
	repo := postgres.NewAIRepository(db.DB)
	lessonRepo := postgres.NewLessonRepository(db.DB)

	// Get all lessons
	lessons, err := lessonRepo.ListByCourse(courseID)
	if err != nil {
		log.Printf("[AI Processing] Failed to get lessons: %v", err)
		errMsg := err.Error()
		repo.UpdateProcessingStatus(ctx, statusID, "failed", 0, 0, &errMsg)
		return
	}
	log.Printf("[AI Processing] Found %d lessons", len(lessons))

	// Clear existing embeddings
	repo.DeleteByCourse(ctx, courseID)

	// Setup processor and embedder based on provider
	proc := processor.NewContentProcessor()
	embConfig := embeddings.EmbedderConfig{
		APIKey: apiKey,
	}

	// Create appropriate embedder with embedding model from settings
	// Note: embedding models are different from chat models!
	var embedder embeddings.Embedder
	embeddingModel := getSettingValue("ai_embedding_model", "")
	
	switch provider {
	case "gemini":
		if embeddingModel == "" {
			embeddingModel = "gemini-embedding-001" // Default Gemini embedding model
		}
		embConfig.Model = embeddingModel
		embedder = embeddings.NewGeminiEmbedder(embConfig)
	default: // openai
		if embeddingModel == "" {
			embeddingModel = "text-embedding-ada-002" // Default OpenAI embedding model
		}
		embConfig.Model = embeddingModel
		embedder = embeddings.NewOpenAIEmbedder(embConfig)
	}

	totalChunks := 0
	processedChunks := 0
	var allChunks []rag.ChunkData

	// Process each lesson
	for _, lesson := range lessons {
		var chunks []rag.Chunk
		log.Printf("[AI Processing] Processing lesson %s, ContentType: %s", lesson.ID, lesson.ContentType)

		switch lesson.ContentType {
		case "text":
			if lesson.Content != nil && *lesson.Content != "" {
				log.Printf("[AI Processing] Lesson %s has text content, length: %d", lesson.ID, len(*lesson.Content))
				chunks = proc.ProcessHTML(*lesson.Content)
			} else {
				log.Printf("[AI Processing] Lesson %s has empty text content", lesson.ID)
			}
		case "pdf":
			if lesson.VideoURL != nil && *lesson.VideoURL != "" {
				pdfPath := *lesson.VideoURL
				log.Printf("[AI Processing] Lesson %s has PDF: %s", lesson.ID, pdfPath)
				
				// Check if it's a local path or URL
				var pdfChunks []rag.Chunk
				var err error
				if strings.HasPrefix(pdfPath, "http://") || strings.HasPrefix(pdfPath, "https://") {
					pdfChunks, err = proc.ProcessPDFFromURL(ctx, pdfPath)
				} else {
					// Local path - prepend /app for container or current dir for local
					// In Docker, files are at /app/uploads/...
					if strings.HasPrefix(pdfPath, "/uploads/") {
						pdfPath = "/app" + pdfPath
					}
					log.Printf("[AI Processing] Processing local PDF: %s", pdfPath)
					pdfChunks, err = proc.ProcessPDFFromFile(pdfPath)
				}
				
				if err == nil {
					chunks = pdfChunks
				} else {
					log.Printf("[AI Processing] PDF processing error: %v", err)
				}
			}
		case "video":
			// Try to get YouTube transcript
			if lesson.VideoURL != nil && *lesson.VideoURL != "" {
				log.Printf("[AI Processing] Lesson %s has video: %s", lesson.ID, *lesson.VideoURL)
				vidChunks, err := proc.ProcessYouTubeTranscript(ctx, *lesson.VideoURL)
				if err == nil {
					chunks = vidChunks
				} else {
					log.Printf("[AI Processing] Video transcript error: %v", err)
				}
			}
		default:
			log.Printf("[AI Processing] Unknown content type: %s for lesson %s", lesson.ContentType, lesson.ID)
		}

		log.Printf("[AI Processing] Lesson %s generated %d chunks", lesson.ID, len(chunks))

		// Skip if no content
		if len(chunks) == 0 {
			continue
		}

		totalChunks += len(chunks)

		// Generate embeddings for chunks
		log.Printf("[AI Processing] Generating embeddings for %d chunks...", len(chunks))
		for i, chunk := range chunks {
			log.Printf("[AI Processing] Processing chunk %d/%d, text length: %d", i+1, len(chunks), len(chunk.Text))
			
			chunkData := rag.ChunkData{
				CourseID:    courseID,
				LessonID:    lesson.ID,
				ContentType: string(lesson.ContentType),
				ChunkIndex:  chunk.Index,
				ChunkText:   chunk.Text,
			}

			// Generate embedding
			log.Printf("[AI Processing] Calling embedder.Embed()...")
			result, err := embedder.Embed(ctx, chunk.Text)
			if err != nil {
				log.Printf("[AI Processing] Embedding error: %v", err)
				continue
			}
			log.Printf("[AI Processing] Embedding generated with %d dimensions", len(result.Embedding))
			chunkData.Embedding = result.Embedding

			allChunks = append(allChunks, chunkData)
			processedChunks++

			// Update progress periodically
			if processedChunks%10 == 0 {
			repo.UpdateProcessingStatus(ctx, statusID, "processing", totalChunks, processedChunks, nil)
			}
		}
	}

	log.Printf("[AI Processing] Total chunks to save: %d, processed: %d", len(allChunks), processedChunks)

	// Batch insert all chunks
	if len(allChunks) > 0 {
		log.Printf("[AI Processing] Inserting %d chunks into database...", len(allChunks))
		err = repo.InsertBatch(ctx, allChunks)
		if err != nil {
			log.Printf("[AI Processing] InsertBatch error: %v", err)
			errMsg := err.Error()
			repo.UpdateProcessingStatus(ctx, statusID, "failed", totalChunks, processedChunks, &errMsg)
			return
		}
		log.Printf("[AI Processing] Chunks inserted successfully!")
	} else {
		log.Printf("[AI Processing] No chunks to insert")
	}

	// Mark as completed
	log.Printf("[AI Processing] Marking status as completed")
	repo.UpdateProcessingStatus(ctx, statusID, "completed", totalChunks, processedChunks, nil)
	log.Printf("[AI Processing] Processing complete!")
}

// GetProcessingStatus returns content processing status for a course
func GetProcessingStatus(c echo.Context) error {
	courseID := c.Param("id")
	ctx := c.Request().Context()

	repo := postgres.NewAIRepository(db.DB)

	statuses, err := repo.GetProcessingStatus(ctx, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get status",
		})
	}

	embeddingCount, _ := repo.GetEmbeddingCount(ctx, courseID)

	// Determine latest status for frontend
	latestStatus := "pending"
	var totalChunks, processedChunks int
	
	// PRIORITY 1: Check if there's an active processing status
	// This must be checked FIRST because during processing, old embeddings may still exist
	if len(statuses) > 0 {
		latest := statuses[0]
		// If status is "processing", always show it (even if embeddings exist from previous run)
		if latest.Status == "processing" {
			latestStatus = "processing"
			totalChunks = latest.TotalChunks
			processedChunks = latest.ProcessedChunks
		} else if latest.Status == "completed" || latest.Status == "failed" {
			// Use status from table
			latestStatus = latest.Status
			totalChunks = latest.TotalChunks
			processedChunks = latest.ProcessedChunks
		}
	}
	
	// PRIORITY 2: If no active status but embeddings exist, show as completed
	if latestStatus == "pending" && embeddingCount > 0 {
		latestStatus = "completed"
		totalChunks = embeddingCount
		processedChunks = embeddingCount
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":           latestStatus,
		"embedding_count":  embeddingCount,
		"total_chunks":     totalChunks,
		"processed_chunks": processedChunks,
		"statuses":         statuses, // Keep for backwards compatibility
	})
}

// ClearCourseEmbeddings removes all embeddings for a course
func ClearCourseEmbeddings(c echo.Context) error {
	courseID := c.Param("id")
	ctx := c.Request().Context()

	repo := postgres.NewAIRepository(db.DB)

	if err := repo.DeleteByCourse(ctx, courseID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to clear embeddings",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Embeddings cleared",
	})
}
