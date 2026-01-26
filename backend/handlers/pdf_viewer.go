package handlers

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/storage"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

// PDFPagesResponse contains info about PDF pages
type PDFPagesResponse struct {
	TotalPages int    `json:"total_pages"`
	LessonID   string `json:"lesson_id"`
}

// GetPDFPages returns the total number of pages in a PDF
// GET /api/content/:lessonId/pdf/pages
func GetPDFPages(c echo.Context) error {
	lessonID := c.Param("lessonId")
	if lessonID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lesson ID is required"})
	}

	// Check access
	userID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get lesson and verify access
	lesson, err := getLessonWithAccess(lessonID, userID, role)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	// Get PDF content URL
	pdfURL := getPDFURL(lesson)
	if pdfURL == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No PDF content found"})
	}

	// Download PDF and count pages
	pageCount, err := getPDFPageCount(pdfURL)
	if err != nil {
		log.Printf("Failed to get PDF page count: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process PDF"})
	}

	return c.JSON(http.StatusOK, PDFPagesResponse{
		TotalPages: pageCount,
		LessonID:   lessonID,
	})
}

// GetPDFPageImage converts a specific PDF page to image and streams it
// GET /api/content/:lessonId/pdf/page/:pageNum
func GetPDFPageImage(c echo.Context) error {
	lessonID := c.Param("lessonId")
	pageNumStr := c.Param("pageNum")
	
	if lessonID == "" || pageNumStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lesson ID and page number are required"})
	}

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil || pageNum < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid page number"})
	}

	// Check access
	userID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get user email for watermark
	userEmail := ""
	userRepo := postgres.NewUserRepository(db.DB)
	user, err := userRepo.GetByID(userID)
	if err == nil && user != nil {
		userEmail = user.Email
	}

	// Get lesson and verify access
	lesson, err := getLessonWithAccess(lessonID, userID, role)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	// Get PDF content URL
	pdfURL := getPDFURL(lesson)
	if pdfURL == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No PDF content found"})
	}

	// Convert page to image
	imgData, err := convertPDFPageToImage(pdfURL, pageNum, userEmail)
	if err != nil {
		log.Printf("Failed to convert PDF page to image: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to render page"})
	}

	// Set headers
	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().Header().Set("Content-Length", strconv.Itoa(len(imgData)))
	c.Response().Header().Set("Cache-Control", "private, no-store")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")

	return c.Blob(http.StatusOK, "image/png", imgData)
}

// Helper function to get lesson with access check
func getLessonWithAccess(lessonID, userID, role string) (*domain.Lesson, error) {
	lessonRepo := postgres.NewLessonRepository(db.DB)
	lesson, err := lessonRepo.GetByID(lessonID)
	if err != nil {
		return nil, fmt.Errorf("lesson not found")
	}

	// Admin can access anything
	if role == "admin" {
		return lesson, nil
	}

	// Check if preview content
	if lesson.IsPreview {
		return lesson, nil
	}

	// Check if instructor
	courseRepo := postgres.NewCourseRepository(db.DB)
	course, err := courseRepo.GetByID(lesson.CourseID)
	if err == nil && course.InstructorID != nil && *course.InstructorID == userID {
		return lesson, nil
	}

	// Check enrollment
	enrollmentRepo := postgres.NewEnrollmentRepository(db.DB)
	enrolled, err := enrollmentRepo.IsEnrolled(userID, lesson.CourseID)
	if err != nil || !enrolled {
		return nil, fmt.Errorf("you must be enrolled in this course to access this content")
	}

	return lesson, nil
}

// Helper to get PDF URL from lesson
func getPDFURL(lesson *domain.Lesson) string {
	if lesson.VideoURL != nil && strings.HasSuffix(strings.ToLower(*lesson.VideoURL), ".pdf") {
		return *lesson.VideoURL
	}
	if lesson.Content != nil && strings.HasSuffix(strings.ToLower(*lesson.Content), ".pdf") {
		return *lesson.Content
	}
	// Also check if VideoURL contains any PDF (sometimes stored there)
	if lesson.VideoURL != nil && *lesson.VideoURL != "" {
		return *lesson.VideoURL
	}
	return ""
}

// downloadPDFToTemp downloads PDF from MinIO to a temp file
func downloadPDFToTemp(objectName string) (string, error) {
	minioStorage := storage.GetStorage()
	if minioStorage == nil {
		return "", fmt.Errorf("storage not configured")
	}

	ctx := context.Background()
	
	// Handle legacy local files
	if strings.HasPrefix(objectName, "/uploads/") {
		// Return local file path
		return "." + objectName, nil
	}

	// Download from MinIO
	obj, err := minioStorage.GetObject(ctx, "", objectName)
	if err != nil {
		return "", fmt.Errorf("failed to get object: %v", err)
	}
	defer obj.Close()

	// Create temp file
	tmpFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	// Copy content
	_, err = io.Copy(tmpFile, obj)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to download PDF: %v", err)
	}

	return tmpFile.Name(), nil
}

// getPDFPageCount returns the number of pages in a PDF
func getPDFPageCount(objectName string) (int, error) {
	// Download PDF to temp
	tmpPath, err := downloadPDFToTemp(objectName)
	if err != nil {
		return 0, err
	}
	
	// Clean up if it's a temp file (not legacy local)
	if !strings.HasPrefix(objectName, "/uploads/") {
		defer os.Remove(tmpPath)
	}

	// Use pdfinfo to get page count
	cmd := exec.Command("pdfinfo", tmpPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("pdfinfo failed: %v, stderr: %s", err, stderr.String())
	}

	// Parse output for "Pages:" line
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Pages:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return strconv.Atoi(parts[1])
			}
		}
	}

	return 0, fmt.Errorf("could not find page count")
}

// convertPDFPageToImage converts a single PDF page to PNG image
func convertPDFPageToImage(objectName string, pageNum int, watermarkEmail string) ([]byte, error) {
	// Download PDF to temp
	tmpPath, err := downloadPDFToTemp(objectName)
	if err != nil {
		return nil, err
	}

	// Clean up temp file if it's not legacy local
	if !strings.HasPrefix(objectName, "/uploads/") {
		defer os.Remove(tmpPath)
	}

	// Create temp output dir
	tmpDir, err := os.MkdirTemp("", "pdfpage-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Convert using pdftoppm (from poppler-utils)
	// -png: output PNG format
	// -f: first page
	// -l: last page
	// -r: resolution (150 DPI is good balance of quality and size)
	outputPrefix := filepath.Join(tmpDir, "page")
	cmd := exec.Command("pdftoppm", 
		"-png", 
		"-f", strconv.Itoa(pageNum), 
		"-l", strconv.Itoa(pageNum),
		"-r", "150",
		tmpPath, 
		outputPrefix,
	)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pdftoppm failed: %v, stderr: %s", err, stderr.String())
	}

	// Find the output file (pdftoppm adds page number suffix)
	files, err := filepath.Glob(filepath.Join(tmpDir, "page-*.png"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no output image found")
	}

	// Read the image
	imgData, err := os.ReadFile(files[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %v", err)
	}

	// If watermark email provided, add watermark
	if watermarkEmail != "" {
		imgData, err = addWatermark(imgData, watermarkEmail)
		if err != nil {
			log.Printf("Warning: failed to add watermark: %v", err)
			// Continue without watermark
		}
	}

	return imgData, nil
}

// addWatermark adds a text watermark to an image
func addWatermark(imgData []byte, email string) ([]byte, error) {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return imgData, err
	}

	// Create new RGBA image
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)

	// Create semi-transparent watermark pattern
	// Using simple approach: add colored rectangles as watermark markers
	watermarkColor := color.RGBA{100, 100, 100, 20} // Very light gray, very transparent

	// Add diagonal watermark pattern across the image
	step := 200 // Distance between watermarks
	for y := 0; y < bounds.Dy(); y += step {
		for x := 0; x < bounds.Dx(); x += step {
			// Draw a small marker at each position
			for dy := 0; dy < 3; dy++ {
				for dx := 0; dx < len(email)*6; dx++ {
					px := x + dx
					py := y + dy
					if px < bounds.Dx() && py < bounds.Dy() {
						rgba.Set(px, py, watermarkColor)
					}
				}
			}
		}
	}

	// Encode back to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
		return imgData, err
	}

	return buf.Bytes(), nil
}
