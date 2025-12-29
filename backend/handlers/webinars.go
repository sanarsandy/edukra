package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
)

var webinarRepo *postgres.WebinarRepository

func initWebinarRepo() {
	if webinarRepo == nil && db.DB != nil {
		webinarRepo = postgres.NewWebinarRepository(db.DB)
	}
}

// ========== ADMIN ENDPOINTS ==========

// CreateWebinar creates a new webinar
// POST /api/admin/webinars
func CreateWebinar(c echo.Context) error {
	initWebinarRepo()

	var req domain.CreateWebinarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.CourseID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "course_id is required"})
	}
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
	}
	if req.ScheduledAt == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "scheduled_at is required"})
	}

	// Parse scheduled_at
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid scheduled_at format. Use ISO 8601 (RFC3339)"})
	}

	// Create webinar
	webinar := &domain.Webinar{
		CourseID:        req.CourseID,
		Title:           req.Title,
		Description:     req.Description,
		ScheduledAt:     scheduledAt,
		DurationMinutes: req.DurationMinutes,
		MeetingURL:      req.MeetingURL,
		MeetingPassword: req.MeetingPassword,
		MaxParticipants: req.MaxParticipants,
		Status:          domain.WebinarStatusUpcoming,
	}

	if req.DurationMinutes == 0 {
		webinar.DurationMinutes = 60 // Default 1 hour
	}

	if req.Status != "" {
		webinar.Status = req.Status
	}

	if err := webinarRepo.Create(webinar); err != nil {
		log.Printf("[Webinar] Failed to create webinar: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create webinar"})
	}

	log.Printf("[Webinar] Created webinar %s for course %s", webinar.ID, webinar.CourseID)
	return c.JSON(http.StatusCreated, webinar)
}

// ListWebinars lists all webinars
// GET /api/admin/webinars
func ListWebinars(c echo.Context) error {
	initWebinarRepo()

	// Parse pagination
	limit := 20
	offset := 0
	if l := c.QueryParam("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := c.QueryParam("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	webinars, total, err := webinarRepo.List(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch webinars"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"webinars": webinars,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetWebinar gets a single webinar by ID
// GET /api/admin/webinars/:id
func GetWebinar(c echo.Context) error {
	initWebinarRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Webinar ID is required"})
	}

	webinar, err := webinarRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch webinar"})
	}
	if webinar == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Webinar not found"})
	}

	return c.JSON(http.StatusOK, webinar)
}

// UpdateWebinar updates a webinar
// PUT /api/admin/webinars/:id
func UpdateWebinar(c echo.Context) error {
	initWebinarRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Webinar ID is required"})
	}

	// Get existing webinar
	webinar, err := webinarRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch webinar"})
	}
	if webinar == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Webinar not found"})
	}

	var req domain.UpdateWebinarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Apply updates
	if req.Title != nil {
		webinar.Title = *req.Title
	}
	if req.Description != nil {
		webinar.Description = req.Description
	}
	if req.ScheduledAt != nil {
		scheduledAt, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid scheduled_at format"})
		}
		webinar.ScheduledAt = scheduledAt
	}
	if req.DurationMinutes != nil {
		webinar.DurationMinutes = *req.DurationMinutes
	}
	if req.MeetingURL != nil {
		webinar.MeetingURL = req.MeetingURL
	}
	if req.MeetingPassword != nil {
		webinar.MeetingPassword = req.MeetingPassword
	}
	if req.MaxParticipants != nil {
		webinar.MaxParticipants = req.MaxParticipants
	}
	if req.Status != nil {
		webinar.Status = *req.Status
	}
	if req.RecordingURL != nil {
		webinar.RecordingURL = req.RecordingURL
	}

	if err := webinarRepo.Update(webinar); err != nil {
		log.Printf("[Webinar] Failed to update webinar: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update webinar"})
	}

	return c.JSON(http.StatusOK, webinar)
}

// DeleteWebinar deletes a webinar
// DELETE /api/admin/webinars/:id
func DeleteWebinar(c echo.Context) error {
	initWebinarRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Webinar ID is required"})
	}

	if err := webinarRepo.Delete(id); err != nil {
		log.Printf("[Webinar] Failed to delete webinar: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete webinar"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Webinar deleted successfully"})
}

// GetWebinarRegistrations gets all registrations for a webinar
// GET /api/admin/webinars/:id/registrations
func GetWebinarRegistrations(c echo.Context) error {
	initWebinarRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Webinar ID is required"})
	}

	registrations, err := webinarRepo.GetRegistrations(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch registrations"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"registrations": registrations,
		"total":         len(registrations),
	})
}

// GetWebinarsByCourse gets all webinars for a course
// GET /api/admin/courses/:id/webinars
func GetWebinarsByCourse(c echo.Context) error {
	initWebinarRepo()

	courseID := c.Param("id")
	if courseID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Course ID is required"})
	}

	webinars, err := webinarRepo.GetByCourseID(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch webinars"})
	}

	return c.JSON(http.StatusOK, webinars)
}

// MarkWebinarAttendance marks user attendance for a webinar
// POST /api/admin/webinars/:id/attendance/:user_id
func MarkWebinarAttendance(c echo.Context) error {
	initWebinarRepo()

	webinarID := c.Param("id")
	userID := c.Param("user_id")

	if webinarID == "" || userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Webinar ID and User ID are required"})
	}

	if err := webinarRepo.MarkAttendance(webinarID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to mark attendance"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Attendance marked successfully"})
}

// ========== STUDENT ENDPOINTS ==========

// GetMyWebinars gets webinars the current user is registered for
// GET /api/my/webinars
func GetMyWebinars(c echo.Context) error {
	initWebinarRepo()

	userID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Check if include completed webinars
	includeCompleted := c.QueryParam("include_completed") == "true"

	webinars, err := webinarRepo.GetUserWebinars(userID, includeCompleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch webinars"})
	}

	return c.JSON(http.StatusOK, webinars)
}

// GetPublicWebinar gets webinar details (public accessible)
// GET /api/webinars/:id
func GetPublicWebinar(c echo.Context) error {
	initWebinarRepo()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Webinar ID is required"})
	}

	webinar, err := webinarRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch webinar"})
	}
	if webinar == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Webinar not found"})
	}

	// Hide sensitive info for public endpoint
	webinar.MeetingPassword = nil

	return c.JSON(http.StatusOK, webinar)
}

// GetCourseWebinars gets upcoming webinars for a course (public)
// GET /api/courses/:id/webinars
func GetCourseWebinars(c echo.Context) error {
	initWebinarRepo()

	courseID := c.Param("id")
	if courseID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Course ID is required"})
	}

	webinars, err := webinarRepo.GetUpcomingByCourse(courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch webinars"})
	}

	// Hide sensitive info
	for _, w := range webinars {
		w.MeetingPassword = nil
		w.MeetingURL = nil // Only show after registration
	}

	return c.JSON(http.StatusOK, webinars)
}
