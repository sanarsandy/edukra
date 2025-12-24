package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
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
}

// UploadFile handles file upload
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

	// Create uploads directory if not exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
	}

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	safeFilename := sanitizeFilename(file.Filename)
	newFilename := fmt.Sprintf("%d_%s", timestamp, safeFilename)
	filePath := filepath.Join(uploadDir, newFilename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	// Return file URL
	return c.JSON(http.StatusOK, UploadResponse{
		URL:      fmt.Sprintf("/uploads/%s", newFilename),
		Filename: file.Filename,
		Size:     file.Size,
		Type:     fileType,
	})
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
