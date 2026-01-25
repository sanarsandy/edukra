package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/storage"
)

// Allowed file extensions by type
var allowedExtensions = map[string][]string{
	"video":    {".mp4", ".webm", ".mov", ".avi", ".mkv"},
	"document": {".pdf", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx"},
	"archive":  {".zip", ".rar", ".7z"},
	"image":    {".jpg", ".jpeg", ".png", ".gif", ".webp"},
}

// Max file sizes (in bytes)
const (
	MaxVideoSize    = 500 * 1024 * 1024 // 500MB
	MaxDocumentSize = 50 * 1024 * 1024  // 50MB
	MaxImageSize    = 10 * 1024 * 1024  // 10MB
)

// UploadResponse represents the response after successful upload
type UploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
	ObjectKey string `json:"object_key,omitempty"` // MinIO object key
}

// UploadFile handles file upload - uploads to MinIO if configured, otherwise local
func UploadFile(c echo.Context) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	
	// Determine file type and validate
	fileType := getFileType(ext)
	if fileType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File type not allowed. Allowed: video (.mp4, .webm, .mov), document (.pdf, .doc, .docx), image (.jpg, .png)",
		})
	}

	// Check file size
	maxSize := getMaxSize(fileType)
	if file.Size > maxSize {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("File too large. Max size for %s: %dMB", fileType, maxSize/(1024*1024)),
		})
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read file"})
	}
	defer src.Close()

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	safeFilename := sanitizeFilename(file.Filename)
	newFilename := fmt.Sprintf("%d_%s", timestamp, safeFilename)

	// Try to upload to MinIO if configured
	minioStorage := storage.GetStorage()
	if minioStorage != nil {
		return uploadToMinio(c, minioStorage, src, newFilename, file.Size, fileType, file.Filename)
	}

	// Fallback to local storage (legacy behavior)
	return uploadToLocal(c, src, newFilename, file.Size, fileType, file.Filename)
}

// uploadToMinio handles upload to MinIO storage
func uploadToMinio(c echo.Context, s *storage.MinioStorage, src io.Reader, objectName string, size int64, fileType, originalFilename string) error {
	ctx := context.Background()
	
	// Determine content type
	contentType := getContentType(objectName)
	
	// Upload to MinIO
	err := s.Upload(ctx, "", objectName, src, size, contentType)
	if err != nil {
		// Log detailed error for debugging
		fmt.Printf("[UPLOAD ERROR] Failed to upload %s: %v\n", objectName, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to upload file to storage: %v", err),
		})
	}

	// Return the object key (not a presigned URL - that will be generated on access)
	return c.JSON(http.StatusOK, UploadResponse{
		URL:       objectName, // Store object key, not URL
		Filename:  originalFilename,
		Size:      size,
		Type:      fileType,
		ObjectKey: objectName,
	})
}

// uploadToLocal handles upload to local filesystem (legacy - deprecated)
func uploadToLocal(c echo.Context, src io.Reader, newFilename string, size int64, fileType, originalFilename string) error {
	// Local storage is deprecated - MinIO must be configured
	return c.JSON(http.StatusServiceUnavailable, map[string]string{
		"error": "Local storage is deprecated. Please configure MinIO for file storage.",
	})
}

// getContentType returns MIME type based on file extension
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".mkv":
		return "video/x-matroska"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/x-rar-compressed"
	case ".7z":
		return "application/x-7z-compressed"
	default:
		return "application/octet-stream"
	}
}

// getFileType determines the file type from extension
func getFileType(ext string) string {
	for fileType, extensions := range allowedExtensions {
		for _, allowedExt := range extensions {
			if ext == allowedExt {
				return fileType
			}
		}
	}
	return ""
}

// getMaxSize returns max allowed size for file type
func getMaxSize(fileType string) int64 {
	switch fileType {
	case "video":
		return MaxVideoSize
	case "document", "archive":
		return MaxDocumentSize
	case "image":
		return MaxImageSize
	default:
		return MaxDocumentSize
	}
}

// sanitizeFilename removes unsafe characters from filename
func sanitizeFilename(filename string) string {
	// Keep only alphanumeric, dots, underscores, and hyphens
	safe := ""
	for _, r := range filename {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			safe += string(r)
		} else if r == ' ' {
			safe += "_"
		}
	}
	return safe
}
