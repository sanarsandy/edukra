package handlers

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/storage"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

// VideoEncryptionKey represents an encryption key record from the database
type VideoEncryptionKey struct {
	ID            string `db:"id" json:"id"`
	LessonID      string `db:"lesson_id" json:"lesson_id"`
	EncryptionKey []byte `db:"encryption_key" json:"-"`
	IV            []byte `db:"iv" json:"-"`
	HLSPath       string `db:"hls_path" json:"hls_path"`
	Status        string `db:"status" json:"status"`
}

// GetHLSManifest serves the HLS manifest (.m3u8) for a lesson
// GET /api/content/:lessonId/hls/manifest
func GetHLSManifest(c echo.Context) error {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lesson ID is required"})
	}

	// Verify user access
	userID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Check enrollment/admin access (reuse logic from secure_content.go)
	if err := verifyContentAccess(userID, role, lessonID); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	// Get encryption key record
	var encKey VideoEncryptionKey
	err = db.DB.Get(&encKey, `
		SELECT id, lesson_id, hls_path, status 
		FROM video_encryption_keys 
		WHERE lesson_id = $1
	`, lessonID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "HLS content not available for this lesson"})
	}
	if err != nil {
		log.Printf("Error fetching HLS key: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch HLS info"})
	}

	if encKey.Status != "ready" {
		return c.JSON(http.StatusAccepted, map[string]string{
			"status":  encKey.Status,
			"message": "Video is still being processed",
		})
	}

	// Get manifest from MinIO
	minioStorage := storage.GetStorage()
	if minioStorage == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Storage not available"})
	}

	manifestPath := encKey.HLSPath + "/playlist.m3u8"
	obj, err := minioStorage.GetObject(context.Background(), "", manifestPath)
	if err != nil {
		log.Printf("Failed to get manifest %s: %v", manifestPath, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Manifest not found"})
	}
	defer obj.Close()

	// Read manifest content
	content, err := io.ReadAll(obj)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read manifest"})
	}

	// Rewrite key URI to point to our authenticated endpoint
	keyURI := fmt.Sprintf("/api/content/%s/hls/key", lessonID)
	modifiedContent := strings.Replace(string(content), 
		`URI="key.bin"`, 
		fmt.Sprintf(`URI="%s"`, keyURI), 
		-1)

	// Set headers
	c.Response().Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")

	return c.String(http.StatusOK, modifiedContent)
}

// GetHLSSegment serves an HLS segment (.ts file)
// GET /api/content/:lessonId/hls/segment/:filename
func GetHLSSegment(c echo.Context) error {
	lessonID := c.Param("lessonId")
	filename := c.Param("filename")

	if lessonID == "" || filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing parameters"})
	}

	// Security: Validate filename to prevent path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid filename"})
	}

	// Only allow .ts files
	if !strings.HasSuffix(filename, ".ts") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid segment type"})
	}

	// Verify user access
	userID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	if err := verifyContentAccess(userID, role, lessonID); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	// Get HLS path
	var hlsPath string
	err = db.DB.Get(&hlsPath, `
		SELECT hls_path FROM video_encryption_keys 
		WHERE lesson_id = $1 AND status = 'ready'
	`, lessonID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "HLS content not found"})
	}

	// Get segment from MinIO
	minioStorage := storage.GetStorage()
	if minioStorage == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Storage not available"})
	}

	segmentPath := filepath.Join(hlsPath, filename)
	obj, err := minioStorage.GetObject(context.Background(), "", segmentPath)
	if err != nil {
		log.Printf("Failed to get segment %s: %v", segmentPath, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Segment not found"})
	}
	defer obj.Close()

	info, err := obj.Stat()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Segment not found"})
	}

	// Set headers
	c.Response().Header().Set("Content-Type", "video/MP2T")
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", info.Size))
	c.Response().Header().Set("Cache-Control", "max-age=86400") // Cache segments for 1 day
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")

	return c.Stream(http.StatusOK, "video/MP2T", obj)
}

// GetHLSKey serves the decryption key for HLS content
// GET /api/content/:lessonId/hls/key
func GetHLSKey(c echo.Context) error {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lesson ID is required"})
	}

	// Verify user access - this MUST be authenticated
	userID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	if err := verifyContentAccess(userID, role, lessonID); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	// Get encryption key from database
	var encKey VideoEncryptionKey
	err = db.DB.Get(&encKey, `
		SELECT encryption_key FROM video_encryption_keys 
		WHERE lesson_id = $1 AND status = 'ready'
	`, lessonID)
	if err != nil {
		log.Printf("Failed to get encryption key for lesson %s: %v", lessonID, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Key not found"})
	}

	// Security headers - prevent caching of key
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")

	// Return the raw key bytes
	return c.Blob(http.StatusOK, "application/octet-stream", encKey.EncryptionKey)
}

// GetHLSStatus returns the HLS processing status for a lesson
// GET /api/content/:lessonId/hls/status
func GetHLSStatus(c echo.Context) error {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lesson ID is required"})
	}

	var status struct {
		Status  string `db:"status" json:"status"`
		HLSPath string `db:"hls_path" json:"hls_path"`
	}

	err := db.DB.Get(&status, `
		SELECT status, COALESCE(hls_path, '') as hls_path
		FROM video_encryption_keys 
		WHERE lesson_id = $1
	`, lessonID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "not_started",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get status"})
	}

	return c.JSON(http.StatusOK, status)
}

// verifyContentAccess checks if a user can access lesson content
func verifyContentAccess(userID, role, lessonID string) error {
	if role == "admin" {
		return nil // Admin can access everything
	}

	// Get lesson info
	var lesson struct {
		CourseID     string  `db:"course_id"`
		IsPreview    bool    `db:"is_preview"`
		InstructorID *string `db:"instructor_id"`
	}

	err := db.DB.Get(&lesson, `
		SELECT l.course_id, l.is_preview, c.instructor_id
		FROM lessons l
		JOIN courses c ON c.id = l.course_id
		WHERE l.id = $1
	`, lessonID)
	if err != nil {
		return fmt.Errorf("lesson not found")
	}

	// Preview content is accessible to all
	if lesson.IsPreview {
		return nil
	}

	// Check if user is the course instructor
	if lesson.InstructorID != nil && *lesson.InstructorID == userID {
		return nil
	}

	// Check enrollment
	var enrolled bool
	err = db.DB.Get(&enrolled, `
		SELECT EXISTS(
			SELECT 1 FROM enrollments 
			WHERE user_id = $1 AND course_id = $2 AND status = 'active'
		)
	`, userID, lesson.CourseID)
	if err != nil || !enrolled {
		return fmt.Errorf("you must be enrolled in this course")
	}

	return nil
}

// ProcessVideoToHLS is called after video upload to start HLS conversion
// This should be called asynchronously (in a goroutine or job queue)
func ProcessVideoToHLS(lessonID, videoPath string) error {
	log.Printf("Starting HLS processing for lesson %s, video: %s", lessonID, videoPath)

	// Insert or update processing status
	_, err := db.DB.Exec(`
		INSERT INTO video_encryption_keys (lesson_id, encryption_key, iv, status)
		VALUES ($1, decode($2, 'hex'), decode($3, 'hex'), 'processing')
		ON CONFLICT (lesson_id) 
		DO UPDATE SET status = 'processing', updated_at = NOW()
	`, lessonID, hex.EncodeToString(make([]byte, 16)), hex.EncodeToString(make([]byte, 16)))
	if err != nil {
		return fmt.Errorf("failed to update processing status: %w", err)
	}

	// Actual processing would happen here
	// For now, return - the actual implementation would:
	// 1. Download video from MinIO to temp
	// 2. Run FFmpeg to create HLS segments
	// 3. Upload segments back to MinIO
	// 4. Update database with ready status

	return nil
}
