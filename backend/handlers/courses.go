package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var courseRepo *postgres.CourseRepository
var lessonRepo *postgres.LessonRepository

func init() {
	// Repositories will be initialized after DB connection
}

func initCourseRepos() {
	if courseRepo == nil && db.DB != nil {
		courseRepo = postgres.NewCourseRepository(db.DB)
		lessonRepo = postgres.NewLessonRepository(db.DB)
	}
}

// ListCourses returns published courses for browsing
func ListCourses(c echo.Context) error {
	initCourseRepos()
	
	tenantID := c.QueryParam("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}
	
	courses, err := courseRepo.ListPublished(tenantID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch courses"})
	}
	
	// Get total count
	total, _ := courseRepo.CountByTenant(tenantID)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"courses": courses,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// AdminListCourses returns all courses for admin management
func AdminListCourses(c echo.Context) error {
	initCourseRepos()
	
	tenantID := c.QueryParam("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}
	
	courses, err := courseRepo.ListByTenant(tenantID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch courses"})
	}
	
	// Get total count
	total, _ := courseRepo.CountByTenant(tenantID)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"courses": courses,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetCourse returns a single course with its lessons
func GetCourse(c echo.Context) error {
	initCourseRepos()
	
	id := c.Param("id")
	course, err := courseRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
	}
	
	if course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}
	
	// Fetch lessons for the course
	lessons, err := lessonRepo.ListByCourse(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch lessons"})
	}
	
	course.Lessons = lessons
	
	return c.JSON(http.StatusOK, course)
}

// CreateCourse creates a new course (admin/instructor only)
func CreateCourse(c echo.Context) error {
	initCourseRepos()
	
	var req domain.CreateCourseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	// Validate required fields
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title is required"})
	}
	
	var instructorID *string
	if req.InstructorID != "" {
		instructorID = &req.InstructorID
	}

	var categoryID *string
	if req.CategoryID != "" {
		categoryID = &req.CategoryID
	}

	course := &domain.Course{
		TenantID:     req.TenantID,
		InstructorID: instructorID,
		CategoryID:   categoryID,
		Title:        req.Title,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailURL,
		Price:        req.Price,
		Currency:     req.Currency,
		LessonsCount: req.LessonsCount,
		Duration:     req.Duration,
		IsPublished:  req.IsPublished,
		IsFeatured:   false,
	}
	
	// if course.TenantID == "" {
	// 	course.TenantID = "default"
	// }
	if course.Currency == "" {
		course.Currency = "IDR"
	}
	
	err := courseRepo.Create(course)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create course: " + err.Error()})
	}
	
	return c.JSON(http.StatusCreated, course)
}

// UpdateCourse updates an existing course
func UpdateCourse(c echo.Context) error {
	initCourseRepos()
	
	id := c.Param("id")
	
	// Fetch existing course
	course, err := courseRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}
	
	var req domain.UpdateCourseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	// Apply updates
	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.Description != nil {
		course.Description = *req.Description
	}
	if req.ThumbnailURL != nil {
		course.ThumbnailURL = req.ThumbnailURL
	}
	if req.Price != nil {
		course.Price = *req.Price
	}
	if req.LessonsCount != nil {
		course.LessonsCount = *req.LessonsCount
	}
	if req.Duration != nil {
		course.Duration = *req.Duration
	}
	if req.IsPublished != nil {
		course.IsPublished = *req.IsPublished
	}
	if req.IsFeatured != nil {
		course.IsFeatured = *req.IsFeatured
	}
	if req.CategoryID != nil {
		course.CategoryID = req.CategoryID
	}
	
	err = courseRepo.Update(course)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update course"})
	}
	
	return c.JSON(http.StatusOK, course)
}

// DeleteCourse deletes a course
func DeleteCourse(c echo.Context) error {
	initCourseRepos()
	
	id := c.Param("id")
	
	// Check if course exists
	course, err := courseRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}
	
	err = courseRepo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete course"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Course deleted successfully"})
}

// PublishCourse toggles the published status
func PublishCourse(c echo.Context) error {
	initCourseRepos()
	
	id := c.Param("id")
	
	course, err := courseRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}
	
	course.IsPublished = !course.IsPublished
	err = courseRepo.Update(course)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update course"})
	}
	
	return c.JSON(http.StatusOK, course)
}
