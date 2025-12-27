package handlers

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var enrollmentRepo *postgres.EnrollmentRepository

func initEnrollmentRepos() {
	if enrollmentRepo == nil && db.DB != nil {
		enrollmentRepo = postgres.NewEnrollmentRepository(db.DB)
	}
	initCourseRepos()
}

// getUserIDFromToken extracts user ID from JWT token
func getUserIDFromToken(c echo.Context) string {
	user := c.Get("user")
	if user == nil {
		return ""
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return ""
	}
	
	// Try pointer first (echojwt uses new(jwt.MapClaims) which returns pointer)
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		// Fallback: try as value type
		if mc, ok := token.Claims.(jwt.MapClaims); ok {
			claims = &mc
		} else {
			return ""
		}
	}
	
	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		return ""
	}
	return userID
}

// ListMyEnrollments returns current user's enrolled courses
func ListMyEnrollments(c echo.Context) error {
	initEnrollmentRepos()
	
	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}
	
	enrollments, err := enrollmentRepo.ListByUser(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch enrollments"})
	}
	
	// Enrich with course data
	type EnrollmentWithCourse struct {
		*postgres.Enrollment
		Course interface{} `json:"course,omitempty"`
	}
	
	var enriched []EnrollmentWithCourse
	for _, e := range enrollments {
		ec := EnrollmentWithCourse{Enrollment: e}
		if courseRepo != nil {
			course, _ := courseRepo.GetByID(e.CourseID)
			ec.Course = course
		}
		enriched = append(enriched, ec)
	}
	
	total, _ := enrollmentRepo.CountByUser(userID)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"enrollments": enriched,
		"total":       total,
		"limit":       limit,
		"offset":      offset,
	})
}

// GetEnrollment returns a single enrollment
func GetEnrollment(c echo.Context) error {
	initEnrollmentRepos()
	
	id := c.Param("id")
	enrollment, err := enrollmentRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch enrollment"})
	}
	
	if enrollment == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Enrollment not found"})
	}
	
	return c.JSON(http.StatusOK, enrollment)
}

// EnrollInCourse creates a new enrollment
func EnrollInCourse(c echo.Context) error {
	initEnrollmentRepos()
	
	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	var req struct {
		CourseID      string  `json:"course_id"`
		TransactionID *string `json:"transaction_id,omitempty"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.CourseID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Course ID is required"})
	}
	
	// Check if already enrolled
	existing, _ := enrollmentRepo.GetByUserAndCourse(userID, req.CourseID)
	if existing != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Already enrolled in this course"})
	}
	
	// Check if course exists
	course, err := courseRepo.GetByID(req.CourseID)
	if err != nil || course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Course not found"})
	}
	
	enrollment := &postgres.Enrollment{
		UserID:        userID,
		CourseID:      req.CourseID,
		TransactionID: req.TransactionID,
	}
	
	err = enrollmentRepo.Create(enrollment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create enrollment: " + err.Error()})
	}
	
	return c.JSON(http.StatusCreated, enrollment)
}

// UpdateEnrollmentProgress updates progress percentage
func UpdateEnrollmentProgress(c echo.Context) error {
	initEnrollmentRepos()
	
	id := c.Param("id")
	
	var req struct {
		Progress int `json:"progress"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.Progress < 0 || req.Progress > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Progress must be between 0 and 100"})
	}
	
	enrollment, err := enrollmentRepo.GetByID(id)
	if err != nil || enrollment == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Enrollment not found"})
	}
	
	if req.Progress >= 100 {
		err = enrollmentRepo.MarkCompleted(id)
	} else {
		err = enrollmentRepo.UpdateProgress(id, req.Progress)
	}
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update progress"})
	}
	
	// Return updated enrollment
	enrollment, _ = enrollmentRepo.GetByID(id)
	return c.JSON(http.StatusOK, enrollment)
}

// CheckEnrollment checks if user is enrolled in a course
func CheckEnrollment(c echo.Context) error {
	initEnrollmentRepos()
	
	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	courseID := c.Param("id")
	
	enrollment, err := enrollmentRepo.GetByUserAndCourse(userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check enrollment"})
	}
	
	if enrollment == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"enrolled": false,
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"enrolled":   true,
		"enrollment": enrollment,
	})
}
