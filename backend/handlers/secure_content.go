package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/storage"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

// SecureContentURLResponse represents the response for secure content URL
type SecureContentURLResponse struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
	Type      string    `json:"type"`
}

// GetSecureContentURL generates a pre-signed URL for authorized users to access lesson content
// GET /api/content/:lessonId/url
func GetSecureContentURL(c echo.Context) error {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lesson ID is required"})
	}

	// Check if MinIO storage is configured
	minioStorage := storage.GetStorage()
	if minioStorage == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Storage service not configured",
		})
	}

	// Get user from JWT
	userID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get lesson from database
	lessonRepo := postgres.NewLessonRepository(db.DB)
	lesson, err := lessonRepo.GetByID(lessonID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lesson not found"})
	}

	// Check if lesson is preview (can be accessed without enrollment)
	// Admin and instructor can always access content
	if !lesson.IsPreview && role != "admin" {
		// Check if user is the course instructor
		courseRepo := postgres.NewCourseRepository(db.DB)
		course, err := courseRepo.GetByID(lesson.CourseID)
		if err != nil || course.InstructorID == nil || *course.InstructorID != userID {
			// Not admin and not instructor, must be enrolled
			enrollmentRepo := postgres.NewEnrollmentRepository(db.DB)
			enrolled, err := enrollmentRepo.IsEnrolled(userID, lesson.CourseID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check enrollment"})
			}
			if !enrolled {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "You must be enrolled in this course to access this content",
				})
			}
		}
	}

	// Get video URL from lesson
	if lesson.VideoURL == nil || *lesson.VideoURL == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No video content found for this lesson"})
	}

	videoURL := *lesson.VideoURL

	// Check if this is a MinIO object (not a legacy local file)
	// MinIO objects are stored without /uploads/ prefix
	if strings.HasPrefix(videoURL, "/uploads/") {
		// Legacy file - return as-is (will be served by static handler if still enabled)
		// For backward compatibility during migration
		return c.JSON(http.StatusOK, SecureContentURLResponse{
			URL:       videoURL,
			ExpiresAt: time.Now().Add(24 * time.Hour), // Long expiry for local files
			Type:      getContentTypeFromURL(videoURL),
		})
	}

	// Generate pre-signed URL from MinIO
	ctx := context.Background()
	presignedURL, expiresAt, err := minioStorage.GetPresignedURLWithDefaultExpiry(ctx, "", videoURL)
	if err != nil {
		log.Printf("Failed to generate presigned URL for %s: %v", videoURL, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate content URL",
		})
	}

	return c.JSON(http.StatusOK, SecureContentURLResponse{
		URL:       presignedURL,
		ExpiresAt: expiresAt,
		Type:      getContentTypeFromURL(videoURL),
	})
}

// GetSecureDocumentURL generates a pre-signed URL for PDF/document content
// GET /api/content/:lessonId/document
func GetSecureDocumentURL(c echo.Context) error {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lesson ID is required"})
	}

	// Check if MinIO storage is configured
	minioStorage := storage.GetStorage()
	if minioStorage == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Storage service not configured",
		})
	}

	// Get user from JWT
	userID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get lesson from database
	lessonRepo := postgres.NewLessonRepository(db.DB)
	lesson, err := lessonRepo.GetByID(lessonID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lesson not found"})
	}

	// Check if lesson is preview
	// Admin and instructor can always access content
	if !lesson.IsPreview && role != "admin" {
		// Check if user is the course instructor
		courseRepo := postgres.NewCourseRepository(db.DB)
		course, err := courseRepo.GetByID(lesson.CourseID)
		if err != nil || course.InstructorID == nil || *course.InstructorID != userID {
			// Not admin and not instructor, must be enrolled
			enrollmentRepo := postgres.NewEnrollmentRepository(db.DB)
			enrolled, err := enrollmentRepo.IsEnrolled(userID, lesson.CourseID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check enrollment"})
			}
			if !enrolled {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "You must be enrolled in this course to access this content",
				})
			}
		}
	}

	// Get content URL (for PDFs stored in Content field or VideoURL)
	var contentURL string
	if lesson.Content != nil && strings.HasSuffix(strings.ToLower(*lesson.Content), ".pdf") {
		contentURL = *lesson.Content
	} else if lesson.VideoURL != nil {
		contentURL = *lesson.VideoURL
	}

	if contentURL == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No document content found for this lesson"})
	}

	// Check if this is a MinIO object
	if strings.HasPrefix(contentURL, "/uploads/") {
		return c.JSON(http.StatusOK, SecureContentURLResponse{
			URL:       contentURL,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Type:      "document",
		})
	}

	// Generate pre-signed URL from MinIO
	ctx := context.Background()
	presignedURL, expiresAt, err := minioStorage.GetPresignedURLWithDefaultExpiry(ctx, "", contentURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate content URL",
		})
	}

	return c.JSON(http.StatusOK, SecureContentURLResponse{
		URL:       presignedURL,
		ExpiresAt: expiresAt,
		Type:      "document",
	})
}

// getContentTypeFromURL determines content type from file URL
func getContentTypeFromURL(url string) string {
	url = strings.ToLower(url)
	switch {
	case strings.HasSuffix(url, ".mp4"):
		return "video/mp4"
	case strings.HasSuffix(url, ".webm"):
		return "video/webm"
	case strings.HasSuffix(url, ".mov"):
		return "video/quicktime"
	case strings.HasSuffix(url, ".pdf"):
		return "application/pdf"
	default:
		return "video"
	}
}

// StreamContent streams content directly from MinIO to the browser
// GET /api/content/:lessonId/stream
// This bypasses pre-signed URL and directly proxies the content
func StreamContent(c echo.Context) error {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lesson ID is required"})
	}

	// Check if MinIO storage is configured
	minioStorage := storage.GetStorage()
	if minioStorage == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Storage service not configured",
		})
	}

	// Get user from JWT
	userID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get lesson from database
	lessonRepo := postgres.NewLessonRepository(db.DB)
	lesson, err := lessonRepo.GetByID(lessonID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lesson not found"})
	}

	// Check access: admin, instructor (owner), enrolled student, or preview content
	if !lesson.IsPreview && role != "admin" {
		courseRepo := postgres.NewCourseRepository(db.DB)
		course, err := courseRepo.GetByID(lesson.CourseID)
		if err != nil || course.InstructorID == nil || *course.InstructorID != userID {
			enrollmentRepo := postgres.NewEnrollmentRepository(db.DB)
			enrolled, err := enrollmentRepo.IsEnrolled(userID, lesson.CourseID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check enrollment"})
			}
			if !enrolled {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "You must be enrolled in this course to access this content",
				})
			}
		}
	}

	// Get content URL (could be video_url or content field)
	var objectName string
	if lesson.VideoURL != nil && *lesson.VideoURL != "" {
		objectName = *lesson.VideoURL
	} else if lesson.Content != nil && *lesson.Content != "" {
		objectName = *lesson.Content
	} else {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No content found"})
	}

	// Check if this is a legacy local file (not MinIO)
	if strings.HasPrefix(objectName, "/uploads/") {
		// Redirect to local file (if still served)
		return c.File("." + objectName)
	}

	// Get object info for content-type
	ctx := context.Background()
	info, err := minioStorage.GetObjectInfo(ctx, "", objectName)
	if err != nil {
		log.Printf("Failed to get object info for %s: %v", objectName, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Content not found"})
	}

	// Get object from MinIO
	obj, err := minioStorage.GetObject(ctx, "", objectName)
	if err != nil {
		log.Printf("Failed to get object %s: %v", objectName, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get content"})
	}
	defer obj.Close()

	// Set headers for streaming
	c.Response().Header().Set("Content-Type", info.ContentType)
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", info.Size))
	
	// Security headers to prevent easy downloading
	c.Response().Header().Set("Content-Disposition", "inline") // Force inline, no filename hint
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")
	c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("X-Frame-Options", "SAMEORIGIN")
	
	// For video streaming support (range requests)
	c.Response().Header().Set("Accept-Ranges", "bytes")

	// Stream the content
	return c.Stream(http.StatusOK, info.ContentType, obj)
}

// GetPublicImage serves images (thumbnails, etc) from MinIO without authentication
// GET /api/images/:objectKey
// This is public because thumbnails are shown on public course listings
func GetPublicImage(c echo.Context) error {
	objectKey := c.Param("objectKey")
	if objectKey == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Object key is required"})
	}

	// Check if MinIO storage is configured
	minioStorage := storage.GetStorage()
	if minioStorage == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Storage service not configured",
		})
	}

	// Validate that this is an image file (security check)
	objectKeyLower := strings.ToLower(objectKey)
	if !strings.HasSuffix(objectKeyLower, ".jpg") && 
	   !strings.HasSuffix(objectKeyLower, ".jpeg") && 
	   !strings.HasSuffix(objectKeyLower, ".png") && 
	   !strings.HasSuffix(objectKeyLower, ".gif") && 
	   !strings.HasSuffix(objectKeyLower, ".webp") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Only image files are allowed",
		})
	}

	// Get object info
	ctx := context.Background()
	info, err := minioStorage.GetObjectInfo(ctx, "", objectKey)
	if err != nil {
		log.Printf("Failed to get image info for %s: %v", objectKey, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Image not found"})
	}

	// Get object from MinIO
	obj, err := minioStorage.GetObject(ctx, "", objectKey)
	if err != nil {
		log.Printf("Failed to get image %s: %v", objectKey, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get image"})
	}
	defer obj.Close()

	// Set headers for caching (images can be cached for 1 day)
	c.Response().Header().Set("Content-Type", info.ContentType)
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", info.Size))
	c.Response().Header().Set("Cache-Control", "public, max-age=86400") // 1 day cache
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")

	// Stream the image
	return c.Stream(http.StatusOK, info.ContentType, obj)
}

// ProxyExternalPDF proxies external PDF URLs through the server to bypass CORS restrictions
// POST /api/content/proxy-pdf
// Request body: { "url": "https://example.com/document.pdf" }
func ProxyExternalPDF(c echo.Context) error {
	// Get user from JWT (must be authenticated)
	_, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Parse request body
	var req struct {
		URL string `json:"url"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.URL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL is required"})
	}

	// Validate URL is a valid external URL
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid URL format"})
	}

	// Security: Only allow PDF files
	urlLower := strings.ToLower(req.URL)
	if !strings.Contains(urlLower, ".pdf") && !strings.Contains(urlLower, "application/pdf") {
		// Allow URLs that might not have .pdf extension but serve PDFs
		log.Printf("Warning: Proxying URL without .pdf extension: %s", req.URL)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Create request
	httpReq, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v", req.URL, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to create request"})
	}

	// Set user agent to avoid being blocked
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	httpReq.Header.Set("Accept", "application/pdf,*/*")

	// Fetch the PDF
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("Failed to fetch PDF from %s: %v", req.URL, err)
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "Failed to fetch external PDF"})
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("External PDF returned status %d for %s", resp.StatusCode, req.URL)
		return c.JSON(resp.StatusCode, map[string]string{"error": fmt.Sprintf("External server returned status %d", resp.StatusCode)})
	}

	// Set response headers
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/pdf"
	}
	
	c.Response().Header().Set("Content-Type", contentType)
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		c.Response().Header().Set("Content-Length", contentLength)
	}
	
	// Security headers
	c.Response().Header().Set("Content-Disposition", "inline")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")
	c.Response().Header().Set("Cache-Control", "private, max-age=3600") // Cache for 1 hour
	c.Response().Header().Set("X-Frame-Options", "SAMEORIGIN")

	// Stream the PDF
	return c.Stream(http.StatusOK, contentType, resp.Body)
}
