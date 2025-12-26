package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/rag"
)

// AIRepository handles AI-related database operations
type AIRepository struct {
	db *sqlx.DB
}

// NewAIRepository creates a new AI repository
func NewAIRepository(db *sqlx.DB) *AIRepository {
	return &AIRepository{db: db}
}

// ================================
// EMBEDDINGS
// ================================

// Search performs vector similarity search for a course
func (r *AIRepository) Search(ctx context.Context, courseID string, embedding []float32, limit int) ([]rag.RetrievedChunk, error) {
	// Build vector string for pgvector
	vectorStr := floatSliceToVector(embedding)

	fmt.Printf("[AIRepo.Search] courseID: %s, embedding dims: %d, limit: %d\n", courseID, len(embedding), limit)
	fmt.Printf("[AIRepo.Search] vectorStr first 50 chars: %.50s...\n", vectorStr)

	query := `
		SELECT 
			ce.id,
			ce.course_id,
			ce.lesson_id,
			COALESCE(l.title, '') as lesson_title,
			ce.content_type,
			ce.chunk_index,
			ce.chunk_text,
			ce.source_reference,
			1 - (ce.embedding <=> $1::vector) as similarity
		FROM course_embeddings ce
		LEFT JOIN lessons l ON ce.lesson_id = l.id
		WHERE ce.course_id = $2
		AND ce.embedding IS NOT NULL
		ORDER BY ce.embedding <=> $1::vector
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, vectorStr, courseID, limit)
	if err != nil {
		fmt.Printf("[AIRepo.Search] Query error: %v\n", err)
		return nil, fmt.Errorf("search query failed: %w", err)
	}
	defer rows.Close()

	var results []rag.RetrievedChunk
	rowCount := 0
	for rows.Next() {
		rowCount++
		var chunk rag.RetrievedChunk
		var sourceRef sql.NullString

		err := rows.Scan(
			&chunk.ID,
			&chunk.CourseID,
			&chunk.LessonID,
			&chunk.LessonTitle,
			&chunk.ContentType,
			&chunk.ChunkIndex,
			&chunk.Text,
			&sourceRef,
			&chunk.SimilarityScore,
		)
		if err != nil {
			fmt.Printf("[AIRepo.Search] Scan error: %v\n", err)
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if sourceRef.Valid {
			chunk.SourceReference = sourceRef.String
		}

		results = append(results, chunk)
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		fmt.Printf("[AIRepo.Search] Rows iteration error: %v\n", err)
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	fmt.Printf("[AIRepo.Search] Processed %d rows, returning %d results\n", rowCount, len(results))
	return results, nil
}

// HybridSearch performs combined vector + keyword search using RRF (Reciprocal Rank Fusion)
// This improves accuracy for exact term searches like NIP, article numbers, codes, etc.
func (r *AIRepository) HybridSearch(ctx context.Context, courseID string, embedding []float32, queryText string, limit int) ([]rag.RetrievedChunk, error) {
	vectorStr := floatSliceToVector(embedding)
	
	fmt.Printf("[AIRepo.HybridSearch] courseID: %s, query: %.30s..., limit: %d\n", courseID, queryText, limit)

	// Hybrid search query using RRF (Reciprocal Rank Fusion)
	// RRF formula: score = sum(1 / (k + rank)) where k=60 is a constant
	query := `
		WITH vector_search AS (
			SELECT 
				id,
				ROW_NUMBER() OVER (ORDER BY embedding <=> $1::vector) as rank,
				1 - (embedding <=> $1::vector) as vector_score
			FROM course_embeddings
			WHERE course_id = $2 AND embedding IS NOT NULL
			LIMIT 20
		),
		keyword_search AS (
			SELECT 
				id,
				ROW_NUMBER() OVER (ORDER BY ts_rank(search_vector, plainto_tsquery('simple', $3)) DESC) as rank,
				ts_rank(search_vector, plainto_tsquery('simple', $3)) as keyword_score
			FROM course_embeddings
			WHERE course_id = $2 
				AND search_vector IS NOT NULL
				AND search_vector @@ plainto_tsquery('simple', $3)
			LIMIT 20
		),
		combined AS (
			SELECT 
				COALESCE(v.id, k.id) as id,
				COALESCE(1.0 / (60 + v.rank), 0) + COALESCE(1.0 / (60 + k.rank), 0) as rrf_score,
				COALESCE(v.vector_score, 0) as vector_score,
				COALESCE(k.keyword_score, 0) as keyword_score
			FROM vector_search v
			FULL OUTER JOIN keyword_search k ON v.id = k.id
		)
		SELECT 
			ce.id,
			ce.course_id,
			ce.lesson_id,
			COALESCE(l.title, '') as lesson_title,
			ce.content_type,
			ce.chunk_index,
			ce.chunk_text,
			ce.source_reference,
			c.rrf_score as similarity
		FROM combined c
		JOIN course_embeddings ce ON ce.id = c.id
		LEFT JOIN lessons l ON ce.lesson_id = l.id
		ORDER BY c.rrf_score DESC
		LIMIT $4
	`

	rows, err := r.db.QueryContext(ctx, query, vectorStr, courseID, queryText, limit)
	if err != nil {
		fmt.Printf("[AIRepo.HybridSearch] Query error: %v\n", err)
		return nil, fmt.Errorf("hybrid search query failed: %w", err)
	}
	defer rows.Close()

	var results []rag.RetrievedChunk
	for rows.Next() {
		var chunk rag.RetrievedChunk
		var sourceRef sql.NullString

		err := rows.Scan(
			&chunk.ID,
			&chunk.CourseID,
			&chunk.LessonID,
			&chunk.LessonTitle,
			&chunk.ContentType,
			&chunk.ChunkIndex,
			&chunk.Text,
			&sourceRef,
			&chunk.SimilarityScore,
		)
		if err != nil {
			fmt.Printf("[AIRepo.HybridSearch] Scan error: %v\n", err)
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if sourceRef.Valid {
			chunk.SourceReference = sourceRef.String
		}

		results = append(results, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	fmt.Printf("[AIRepo.HybridSearch] Returning %d results\n", len(results))
	return results, nil
}

// Insert adds a new embedding
func (r *AIRepository) Insert(ctx context.Context, chunk rag.ChunkData) error {
	vectorStr := floatSliceToVector(chunk.Embedding)

	metadataJSON, _ := json.Marshal(chunk.Metadata)

	query := `
		INSERT INTO course_embeddings (
			course_id, lesson_id, content_type, chunk_index, 
			chunk_text, embedding, source_reference, metadata
		) VALUES ($1, $2, $3, $4, $5, $6::vector, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		chunk.CourseID,
		nullString(chunk.LessonID),
		chunk.ContentType,
		chunk.ChunkIndex,
		chunk.ChunkText,
		vectorStr,
		nullString(chunk.SourceReference),
		metadataJSON,
	)

	return err
}

// InsertBatch adds multiple embeddings
func (r *AIRepository) InsertBatch(ctx context.Context, chunks []rag.ChunkData) error {
	if len(chunks) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO course_embeddings (
			course_id, lesson_id, content_type, chunk_index, 
			chunk_text, embedding, source_reference, metadata
		) VALUES ($1, $2, $3, $4, $5, $6::vector, $7, $8)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, chunk := range chunks {
		vectorStr := floatSliceToVector(chunk.Embedding)
		metadataJSON, _ := json.Marshal(chunk.Metadata)

		// Sanitize text to remove invalid UTF-8 sequences
		cleanText := strings.ToValidUTF8(chunk.ChunkText, "")

		_, err := stmt.ExecContext(ctx,
			chunk.CourseID,
			nullString(chunk.LessonID),
			chunk.ContentType,
			chunk.ChunkIndex,
			cleanText,
			vectorStr,
			nullString(chunk.SourceReference),
			metadataJSON,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteByCourse removes all embeddings for a course
func (r *AIRepository) DeleteByCourse(ctx context.Context, courseID string) error {
	query := `DELETE FROM course_embeddings WHERE course_id = $1`
	_, err := r.db.ExecContext(ctx, query, courseID)
	return err
}

// DeleteByLesson removes all embeddings for a lesson
func (r *AIRepository) DeleteByLesson(ctx context.Context, lessonID string) error {
	query := `DELETE FROM course_embeddings WHERE lesson_id = $1`
	_, err := r.db.ExecContext(ctx, query, lessonID)
	return err
}

// GetEmbeddingCount returns total embeddings for a course
func (r *AIRepository) GetEmbeddingCount(ctx context.Context, courseID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM course_embeddings WHERE course_id = $1`
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&count)
	return count, err
}

// ================================
// PROCESSING STATUS
// ================================

// ProcessingStatus represents content processing status
type ProcessingStatus struct {
	ID              string     `db:"id" json:"id"`
	CourseID        string     `db:"course_id" json:"course_id"`
	LessonID        *string    `db:"lesson_id" json:"lesson_id,omitempty"`
	Status          string     `db:"status" json:"status"`
	ContentType     string     `db:"content_type" json:"content_type"`
	TotalChunks     int        `db:"total_chunks" json:"total_chunks"`
	ProcessedChunks int        `db:"processed_chunks" json:"processed_chunks"`
	ErrorMessage    *string    `db:"error_message" json:"error_message,omitempty"`
	StartedAt       *time.Time `db:"started_at" json:"started_at,omitempty"`
	CompletedAt     *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}

// CreateProcessingStatus creates a new processing status record
func (r *AIRepository) CreateProcessingStatus(ctx context.Context, courseID, lessonID, contentType string) (*ProcessingStatus, error) {
	query := `
		INSERT INTO content_processing_status (course_id, lesson_id, status, content_type)
		VALUES ($1, $2, 'pending', $3)
		RETURNING id, created_at
	`

	var status ProcessingStatus
	status.CourseID = courseID
	if lessonID != "" {
		status.LessonID = &lessonID
	}
	status.Status = "pending"
	status.ContentType = contentType

	err := r.db.QueryRowContext(ctx, query, courseID, nullString(lessonID), contentType).
		Scan(&status.ID, &status.CreatedAt)

	return &status, err
}

// UpdateProcessingStatus updates a processing status record
func (r *AIRepository) UpdateProcessingStatus(ctx context.Context, id, newStatus string, totalChunks, processedChunks int, errorMsg *string) error {
	// Use explicit casting and separate variables to avoid PostgreSQL parameter type inference issues
	query := `
		UPDATE content_processing_status SET
			status = $2::text,
			total_chunks = $3,
			processed_chunks = $4,
			error_message = $5,
			started_at = CASE WHEN status = 'pending' AND $2::text = 'processing' THEN NOW() ELSE started_at END,
			completed_at = CASE WHEN $2::text IN ('completed', 'failed') THEN NOW() ELSE completed_at END
		WHERE id = $1::uuid
	`

	_, err := r.db.ExecContext(ctx, query, id, newStatus, totalChunks, processedChunks, errorMsg)
	return err
}

// GetProcessingStatus gets processing status for a course
func (r *AIRepository) GetProcessingStatus(ctx context.Context, courseID string) ([]ProcessingStatus, error) {
	query := `
		SELECT id, course_id, lesson_id, status, content_type, 
		       total_chunks, processed_chunks, error_message,
		       started_at, completed_at, created_at
		FROM content_processing_status
		WHERE course_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []ProcessingStatus
	for rows.Next() {
		var s ProcessingStatus
		err := rows.Scan(
			&s.ID, &s.CourseID, &s.LessonID, &s.Status, &s.ContentType,
			&s.TotalChunks, &s.ProcessedChunks, &s.ErrorMessage,
			&s.StartedAt, &s.CompletedAt, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}

// ================================
// CHAT SESSIONS
// ================================

// ChatSession represents an AI chat session
type ChatSession struct {
	ID              string    `db:"id" json:"id"`
	UserID          string    `db:"user_id" json:"user_id"`
	CourseID        string    `db:"course_id" json:"course_id"`
	Title           string    `db:"title" json:"title"`
	IsActive        bool      `db:"is_active" json:"is_active"`
	MessageCount    int       `db:"message_count" json:"message_count"`
	TotalTokensUsed int       `db:"total_tokens_used" json:"total_tokens_used"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// ChatMessage represents a message in a chat session
type ChatMessage struct {
	ID         string          `db:"id" json:"id"`
	SessionID  string          `db:"session_id" json:"session_id"`
	Role       string          `db:"role" json:"role"`
	Content    string          `db:"content" json:"content"`
	Sources    json.RawMessage `db:"sources" json:"sources,omitempty"`
	TokensUsed int             `db:"tokens_used" json:"tokens_used"`
	Provider   *string         `db:"provider" json:"provider,omitempty"`
	Model      *string         `db:"model" json:"model,omitempty"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
}

// GetOrCreateActiveSession gets or creates an active chat session
func (r *AIRepository) GetOrCreateActiveSession(ctx context.Context, userID, courseID string) (*ChatSession, error) {
	// Try to get existing active session
	query := `
		SELECT id, user_id, course_id, title, is_active, message_count, 
		       total_tokens_used, created_at, updated_at
		FROM ai_chat_sessions
		WHERE user_id = $1 AND course_id = $2 AND is_active = true
	`

	var session ChatSession
	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(
		&session.ID, &session.UserID, &session.CourseID, &session.Title,
		&session.IsActive, &session.MessageCount, &session.TotalTokensUsed,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new session
		return r.CreateChatSession(ctx, userID, courseID)
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// CreateChatSession creates a new chat session
func (r *AIRepository) CreateChatSession(ctx context.Context, userID, courseID string) (*ChatSession, error) {
	query := `
		INSERT INTO ai_chat_sessions (user_id, course_id, title, is_active)
		VALUES ($1, $2, 'Chat Baru', true)
		RETURNING id, user_id, course_id, title, is_active, message_count, 
		          total_tokens_used, created_at, updated_at
	`

	var session ChatSession
	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(
		&session.ID, &session.UserID, &session.CourseID, &session.Title,
		&session.IsActive, &session.MessageCount, &session.TotalTokensUsed,
		&session.CreatedAt, &session.UpdatedAt,
	)

	return &session, err
}

// AddChatMessage adds a message to a session
func (r *AIRepository) AddChatMessage(ctx context.Context, sessionID, role, content string, sources []rag.Source, tokensUsed int, provider, model string) (*ChatMessage, error) {
	sourcesJSON, _ := json.Marshal(sources)

	query := `
		INSERT INTO ai_chat_messages (session_id, role, content, sources, tokens_used, provider, model)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	var msg ChatMessage
	msg.SessionID = sessionID
	msg.Role = role
	msg.Content = content
	msg.Sources = sourcesJSON
	msg.TokensUsed = tokensUsed

	if provider != "" {
		msg.Provider = &provider
	}
	if model != "" {
		msg.Model = &model
	}

	err := r.db.QueryRowContext(ctx, query,
		sessionID, role, content, sourcesJSON, tokensUsed,
		nullString(provider), nullString(model),
	).Scan(&msg.ID, &msg.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Update session stats
	updateQuery := `
		UPDATE ai_chat_sessions SET
			message_count = message_count + 1,
			total_tokens_used = total_tokens_used + $2,
			updated_at = NOW()
		WHERE id = $1
	`
	r.db.ExecContext(ctx, updateQuery, sessionID, tokensUsed)

	return &msg, nil
}

// GetChatHistory gets messages for a session
func (r *AIRepository) GetChatHistory(ctx context.Context, sessionID string, limit int) ([]ChatMessage, error) {
	query := `
		SELECT id, session_id, role, content, sources, tokens_used, provider, model, created_at
		FROM ai_chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(
			&msg.ID, &msg.SessionID, &msg.Role, &msg.Content,
			&msg.Sources, &msg.TokensUsed, &msg.Provider, &msg.Model, &msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// ClearChatSession deactivates a session
func (r *AIRepository) ClearChatSession(ctx context.Context, sessionID string) error {
	query := `UPDATE ai_chat_sessions SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, sessionID)
	return err
}

// ================================
// USAGE LOGGING
// ================================

// LogUsage logs AI usage for rate limiting and analytics
func (r *AIRepository) LogUsage(ctx context.Context, userID, courseID, actionType, provider, model string, tokensInput, tokensOutput int) error {
	query := `
		INSERT INTO ai_usage_log (user_id, course_id, action_type, provider, model, tokens_input, tokens_output)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		userID, nullString(courseID), actionType, provider, model, tokensInput, tokensOutput,
	)
	return err
}

// GetTodayUsageCount gets chat count for user today
func (r *AIRepository) GetTodayUsageCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM ai_usage_log
		WHERE user_id = $1 
		AND action_type = 'chat'
		AND DATE(created_at) = CURRENT_DATE
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// ================================
// HELPERS
// ================================

// floatSliceToVector converts []float32 to pgvector string format
func floatSliceToVector(v []float32) string {
	if len(v) == 0 {
		return "[]"
	}
	
	parts := make([]string, len(v))
	for i, f := range v {
		// Use strconv for consistent formatting
		parts[i] = strconv.FormatFloat(float64(f), 'f', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// nullString returns nil for empty strings
func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
