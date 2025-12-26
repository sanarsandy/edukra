package rag

import (
	"context"
	"fmt"

	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/embeddings"
)

// RetrievedChunk represents a chunk retrieved from vector search
type RetrievedChunk struct {
	ID             string  `json:"id"`
	CourseID       string  `json:"course_id"`
	LessonID       string  `json:"lesson_id,omitempty"`
	LessonTitle    string  `json:"lesson_title,omitempty"`
	ContentType    string  `json:"content_type"`
	ChunkIndex     int     `json:"chunk_index"`
	Text           string  `json:"text"`
	SimilarityScore float64 `json:"similarity_score"`
	SourceReference string  `json:"source_reference,omitempty"`
}

// Source represents a source reference for citations
type Source struct {
	LessonID    string  `json:"lesson_id"`
	LessonTitle string  `json:"lesson_title"`
	ContentType string  `json:"content_type"`
	Preview     string  `json:"preview"`
	Relevance   float64 `json:"relevance"`
}

// RetrieverConfig holds configuration for the retriever
type RetrieverConfig struct {
	TopK            int     // Number of chunks to retrieve
	MinSimilarity   float64 // Minimum similarity threshold
	MaxContextChars int     // Maximum total context characters
}

// DefaultRetrieverConfig returns sensible defaults
func DefaultRetrieverConfig() RetrieverConfig {
	return RetrieverConfig{
		TopK:            5,
		MinSimilarity:   0.3, // Lowered from 0.7 to allow more relevant chunks
		MaxContextChars: 8000, // ~2000 tokens
	}
}

// VectorStore interface for vector database operations
type VectorStore interface {
	// Search performs similarity search
	Search(ctx context.Context, courseID string, embedding []float32, limit int) ([]RetrievedChunk, error)
	
	// HybridSearch performs combined vector + keyword search using RRF
	HybridSearch(ctx context.Context, courseID string, embedding []float32, queryText string, limit int) ([]RetrievedChunk, error)
	
	// Insert adds a new embedding
	Insert(ctx context.Context, chunk ChunkData) error
	
	// InsertBatch adds multiple embeddings
	InsertBatch(ctx context.Context, chunks []ChunkData) error
	
	// DeleteByCourse removes all embeddings for a course
	DeleteByCourse(ctx context.Context, courseID string) error
	
	// DeleteByLesson removes all embeddings for a lesson
	DeleteByLesson(ctx context.Context, lessonID string) error
}

// ChunkData represents data to be stored in vector DB
type ChunkData struct {
	CourseID        string
	LessonID        string
	ContentType     string
	ChunkIndex      int
	ChunkText       string
	Embedding       []float32
	SourceReference string
	Metadata        map[string]interface{}
}

// Retriever handles RAG retrieval operations
type Retriever struct {
	embedder embeddings.Embedder
	store    VectorStore
	config   RetrieverConfig
}

// NewRetriever creates a new retriever
func NewRetriever(embedder embeddings.Embedder, store VectorStore, config RetrieverConfig) *Retriever {
	return &Retriever{
		embedder: embedder,
		store:    store,
		config:   config,
	}
}

// Retrieve finds relevant chunks for a query using Hybrid Search (Vector + Keyword)
func (r *Retriever) Retrieve(ctx context.Context, courseID, query string) ([]RetrievedChunk, error) {
	// Generate embedding for query
	fmt.Printf("[Retriever] Generating embedding for query: %s\n", query)
	result, err := r.embedder.Embed(ctx, query)
	if err != nil {
		fmt.Printf("[Retriever] Embed error: %v\n", err)
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}
	fmt.Printf("[Retriever] Got embedding with %d dimensions\n", len(result.Embedding))

	// Use Hybrid Search (Vector + Keyword with RRF)
	fmt.Printf("[Retriever] HybridSearch for courseID: %s, limit: %d\n", courseID, r.config.TopK)
	chunks, err := r.store.HybridSearch(ctx, courseID, result.Embedding, query, r.config.TopK)
	if err != nil {
		// Fallback to regular vector search if hybrid fails
		fmt.Printf("[Retriever] HybridSearch error, falling back to vector search: %v\n", err)
		chunks, err = r.store.Search(ctx, courseID, result.Embedding, r.config.TopK)
		if err != nil {
			fmt.Printf("[Retriever] Search error: %v\n", err)
			return nil, fmt.Errorf("failed to search: %w", err)
		}
	}
	fmt.Printf("[Retriever] Search returned %d chunks\n", len(chunks))

	// Filter by minimum similarity (RRF scores are typically smaller, so we adjust threshold)
	var filtered []RetrievedChunk
	minScore := r.config.MinSimilarity * 0.01 // Adjust for RRF scoring (scores are ~0.01-0.03 range)
	for _, chunk := range chunks {
		fmt.Printf("[Retriever] Chunk %s score: %.4f (threshold: %.4f)\n", chunk.ID, chunk.SimilarityScore, minScore)
		if chunk.SimilarityScore >= minScore {
			filtered = append(filtered, chunk)
		}
	}
	fmt.Printf("[Retriever] After filtering: %d chunks\n", len(filtered))

	return filtered, nil
}

// BuildContext creates context string from retrieved chunks
func (r *Retriever) BuildContext(chunks []RetrievedChunk) (string, []Source) {
	if len(chunks) == 0 {
		return "", nil
	}

	var context string
	var sources []Source
	totalChars := 0
	seenLessons := make(map[string]bool)

	for i, chunk := range chunks {
		// Check if we've exceeded max context
		if totalChars+len(chunk.Text) > r.config.MaxContextChars {
			break
		}

		// Add to context
		context += fmt.Sprintf("[Sumber %d: %s]\n%s\n\n", i+1, chunk.LessonTitle, chunk.Text)
		totalChars += len(chunk.Text)

		// Add unique source
		if !seenLessons[chunk.LessonID] {
			sources = append(sources, Source{
				LessonID:    chunk.LessonID,
				LessonTitle: chunk.LessonTitle,
				ContentType: chunk.ContentType,
				Preview:     truncateText(chunk.Text, 100),
				Relevance:   chunk.SimilarityScore,
			})
			seenLessons[chunk.LessonID] = true
		}
	}

	return context, sources
}

// truncateText truncates text to maxLen with ellipsis
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}
