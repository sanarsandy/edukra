package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

func initLessonRepos() {
	if lessonRepo == nil && db.DB != nil {
		lessonRepo = postgres.NewLessonRepository(db.DB)
	}
}

// GetLesson returns a single lesson
func GetLesson(c echo.Context) error {
	initLessonRepos()

	id := c.Param("id")
	lesson, err := lessonRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch lesson"})
	}

	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lesson not found"})
	}

	return c.JSON(http.StatusOK, lesson)
}

// ListLessons returns lessons for a course (supports ?tree=true for nested structure)
func ListLessons(c echo.Context) error {
	initLessonRepos()

	courseID := c.Param("courseId")
	treeMode := c.QueryParam("tree") == "true"

	var lessons []*domain.Lesson
	var err error

	if treeMode {
		lessons, err = lessonRepo.GetTree(courseID)
	} else {
		lessons, err = lessonRepo.ListByCourse(courseID)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch lessons"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lessons": lessons,
		"total":   len(lessons),
	})
}

// GetLessonTree returns lessons for a course as a nested tree structure
func GetLessonTree(c echo.Context) error {
	initLessonRepos()

	courseID := c.Param("courseId")
	lessons, err := lessonRepo.GetTree(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch lessons tree"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lessons": lessons,
		"total":   countLessonsRecursive(lessons),
	})
}

// countLessonsRecursive counts total lessons including nested ones
func countLessonsRecursive(lessons []*domain.Lesson) int {
	count := len(lessons)
	for _, lesson := range lessons {
		if lesson.Children != nil {
			count += countLessonsRecursive(lesson.Children)
		}
	}
	return count
}

// CreateLesson creates a new lesson for a course
func CreateLesson(c echo.Context) error {
	initLessonRepos()

	courseID := c.Param("courseId")

	var req domain.CreateLessonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title is required"})
	}

	// Get current lesson count for order_index (within the same parent)
	count, _ := lessonRepo.CountByParent(courseID, req.ParentID)

	lesson := &domain.Lesson{
		CourseID:      courseID,
		ParentID:      req.ParentID,
		IsContainer:   req.IsContainer,
		Title:         req.Title,
		Description:   req.Description,
		OrderIndex:    count,
		ContentType:   req.ContentType,
		VideoURL:      req.VideoURL,
		Content:       req.Content,
		SecurityLevel: req.SecurityLevel,
		IsPreview:     req.IsPreview,
	}

	// Default content type for containers
	if lesson.IsContainer {
		lesson.ContentType = "" // Containers don't have content type
	} else if lesson.ContentType == "" {
		lesson.ContentType = domain.ContentVideo
	}

	if lesson.SecurityLevel == "" {
		lesson.SecurityLevel = domain.SecuritySignedURL
	}

	err := lessonRepo.Create(lesson)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create lesson: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, lesson)
}

// UpdateLesson updates an existing lesson
func UpdateLesson(c echo.Context) error {
	initLessonRepos()

	id := c.Param("id")

	lesson, err := lessonRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch lesson"})
	}
	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lesson not found"})
	}

	var req domain.UpdateLessonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Apply updates
	if req.ParentID != nil {
		lesson.ParentID = req.ParentID
	}
	if req.IsContainer != nil {
		lesson.IsContainer = *req.IsContainer
	}
	if req.Title != nil {
		lesson.Title = *req.Title
	}
	if req.Description != nil {
		lesson.Description = *req.Description
	}
	if req.OrderIndex != nil {
		lesson.OrderIndex = *req.OrderIndex
	}
	if req.ContentType != nil {
		lesson.ContentType = *req.ContentType
	}
	if req.VideoURL != nil {
		lesson.VideoURL = req.VideoURL
	}
	if req.Content != nil {
		lesson.Content = req.Content
	}
	if req.SecurityLevel != nil {
		lesson.SecurityLevel = *req.SecurityLevel
	}
	if req.IsPreview != nil {
		lesson.IsPreview = *req.IsPreview
	}

	err = lessonRepo.Update(lesson)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update lesson"})
	}

	return c.JSON(http.StatusOK, lesson)
}

// DeleteLesson deletes a lesson (and all children via CASCADE)
func DeleteLesson(c echo.Context) error {
	initLessonRepos()

	id := c.Param("id")

	lesson, err := lessonRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch lesson"})
	}
	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lesson not found"})
	}

	err = lessonRepo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete lesson"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Lesson deleted successfully"})
}

// ReorderLessons updates the order of lessons in a course
func ReorderLessons(c echo.Context) error {
	initLessonRepos()

	courseID := c.Param("courseId")

	var req struct {
		LessonIDs []string `json:"lesson_ids"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := lessonRepo.ReorderLessons(courseID, req.LessonIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reorder lessons"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Lessons reordered successfully"})
}

// MoveLesson moves a lesson to a new parent
func MoveLesson(c echo.Context) error {
	initLessonRepos()

	lessonID := c.Param("id")

	lesson, err := lessonRepo.GetByID(lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch lesson"})
	}
	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lesson not found"})
	}

	var req domain.MoveLessonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Prevent moving a lesson to be its own child (circular reference)
	if req.ParentID != nil && *req.ParentID == lessonID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot move lesson into itself"})
	}

	err = lessonRepo.MoveLesson(lessonID, req.ParentID, req.OrderIndex)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to move lesson"})
	}

	// Return updated lesson
	lesson, _ = lessonRepo.GetByID(lessonID)
	return c.JSON(http.StatusOK, lesson)
}

// ========================================
// SHARED FUNCTIONS (used by admin and instructor handlers)
// ========================================

// createLessonForCourse is a shared function for creating lessons
func createLessonForCourse(c echo.Context, courseID string) error {
	initLessonRepos()

	var req domain.CreateLessonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Judul materi wajib diisi"})
	}

	// Get current lesson count for order_index
	count, _ := lessonRepo.CountByParent(courseID, req.ParentID)

	lesson := &domain.Lesson{
		CourseID:      courseID,
		ParentID:      req.ParentID,
		IsContainer:   req.IsContainer,
		Title:         req.Title,
		Description:   req.Description,
		OrderIndex:    count,
		ContentType:   req.ContentType,
		VideoURL:      req.VideoURL,
		Content:       req.Content,
		SecurityLevel: req.SecurityLevel,
		IsPreview:     req.IsPreview,
	}

	if lesson.IsContainer {
		lesson.ContentType = ""
	} else if lesson.ContentType == "" {
		lesson.ContentType = domain.ContentVideo
	}

	if lesson.SecurityLevel == "" {
		lesson.SecurityLevel = domain.SecuritySignedURL
	}

	err := lessonRepo.Create(lesson)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat materi: " + err.Error()})
	}

	// Update lessons_count in course
	updateCourseLessonsCount(courseID)

	return c.JSON(http.StatusCreated, lesson)
}

// updateLessonByID is a shared function for updating lessons
func updateLessonByID(c echo.Context, lessonID string) error {
	initLessonRepos()

	lesson, err := lessonRepo.GetByID(lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil materi"})
	}
	if lesson == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Materi tidak ditemukan"})
	}

	var req domain.UpdateLessonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	// Apply updates
	if req.ParentID != nil {
		lesson.ParentID = req.ParentID
	}
	if req.IsContainer != nil {
		lesson.IsContainer = *req.IsContainer
	}
	if req.Title != nil {
		lesson.Title = *req.Title
	}
	if req.Description != nil {
		lesson.Description = *req.Description
	}
	if req.OrderIndex != nil {
		lesson.OrderIndex = *req.OrderIndex
	}
	if req.ContentType != nil {
		lesson.ContentType = *req.ContentType
	}
	if req.VideoURL != nil {
		lesson.VideoURL = req.VideoURL
	}
	if req.Content != nil {
		lesson.Content = req.Content
	}
	if req.SecurityLevel != nil {
		lesson.SecurityLevel = *req.SecurityLevel
	}
	if req.IsPreview != nil {
		lesson.IsPreview = *req.IsPreview
	}

	err = lessonRepo.Update(lesson)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal update materi"})
	}

	return c.JSON(http.StatusOK, lesson)
}

// updateCourseLessonsCount updates the lessons_count field in course
func updateCourseLessonsCount(courseID string) {
	var count int
	db.DB.QueryRow(`SELECT COUNT(*) FROM lessons WHERE course_id = $1`, courseID).Scan(&count)
	db.DB.Exec(`UPDATE courses SET lessons_count = $1, updated_at = NOW() WHERE id = $2`, count, courseID)
}
