package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
)

var instructorRepo *postgres.UserRepository

func initInstructorRepos() {
	if instructorRepo == nil && db.DB != nil {
		instructorRepo = postgres.NewUserRepository(db.DB)
	}
}

// InstructorResponse represents an instructor for API response
type InstructorResponse struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	IsActive    bool    `json:"is_active"`
	Specialty   string  `json:"specialty,omitempty"`
	Bio         string  `json:"bio,omitempty"`
	CourseCount int     `json:"course_count"`
	Students    int     `json:"students"`
	Rating      float64 `json:"rating"`
}

func userToInstructor(user *domain.User, courseCount int, studentCount int, rating float64) *InstructorResponse {
	specialty, _ := user.Metadata["specialty"].(string)
	bio, _ := user.Metadata["bio"].(string)

	return &InstructorResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		AvatarURL:   user.AvatarURL,
		IsActive:    user.IsActive,
		Specialty:   specialty,
		Bio:         bio,
		CourseCount: courseCount,
		Students:    studentCount,
		Rating:      rating,
	}
}

// ListInstructors returns all users with role=instructor
func ListInstructors(c echo.Context) error {
	initInstructorRepos()
	
	// Initialize course repo for stats
	courseRepo := postgres.NewCourseRepository(db.DB)
	
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
	
	// Get all users and filter by role
	users, err := instructorRepo.ListByTenant(tenantID, 100, 0)
	if err != nil {
		c.Logger().Errorf("ListInstructors error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch instructors: " + err.Error()})
	}
	
	// Filter instructors and calculate stats
	var instructors []*InstructorResponse
	for _, user := range users {
		if user.Role == "instructor" && user.IsActive {
			// Get course count for this instructor
			courseCount, _ := courseRepo.CountByInstructor(user.ID)
			
			// Get student count for this instructor
			studentCount, _ := courseRepo.GetStudentCountByInstructor(user.ID)
			
			// Get average rating for this instructor
			rating, _ := courseRepo.GetAverageRatingByInstructor(user.ID)
			
			instructors = append(instructors, userToInstructor(user, courseCount, studentCount, rating))
		}
	}
	
	// Apply pagination manually
	start := offset
	end := offset + limit
	if start > len(instructors) {
		start = len(instructors)
	}
	if end > len(instructors) {
		end = len(instructors)
	}
	paginatedInstructors := instructors[start:end]
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"instructors": paginatedInstructors,
		"total":       len(instructors),
	})
}

// GetInstructor returns a single instructor by ID
func GetInstructor(c echo.Context) error {
	initInstructorRepos()
	
	// Initialize course repo for stats
	courseRepo := postgres.NewCourseRepository(db.DB)
	
	id := c.Param("id")
	
	user, err := instructorRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch instructor"})
	}
	
	if user == nil || user.Role != "instructor" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Instructor not found"})
	}
	
	// Get course count and student count
	courseCount, _ := courseRepo.CountByInstructor(id)
	studentCount, _ := courseRepo.GetStudentCountByInstructor(id)
	rating, _ := courseRepo.GetAverageRatingByInstructor(id)
	
	return c.JSON(http.StatusOK, userToInstructor(user, courseCount, studentCount, rating))
}

// CreateInstructor creates a new user with role=instructor
func CreateInstructor(c echo.Context) error {
	initInstructorRepos()
	
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Specialty string `json:"specialty"`
		Bio       string `json:"bio"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}
	if req.FullName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Full name is required"})
	}
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password is required"})
	}
	
	// Check if email already exists
	existingUser, _ := instructorRepo.GetByEmail("default", req.Email)
	if existingUser != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email already registered"})
	}
	
	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	
	// Create user with instructor role
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         "instructor",
		AuthProvider: "email",
		IsActive:     true,
		Metadata: map[string]interface{}{
			"specialty": req.Specialty,
			"bio":       req.Bio,
		},
	}
	
	err = instructorRepo.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create instructor: " + err.Error()})
	}
	
	return c.JSON(http.StatusCreated, userToInstructor(user, 0, 0, 0.0))
}

// UpdateInstructor updates an instructor
func UpdateInstructor(c echo.Context) error {
	initInstructorRepos()
	
	id := c.Param("id")
	
	user, err := instructorRepo.GetByID(id)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Instructor not found"})
	}
	
	if user.Role != "instructor" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User is not an instructor"})
	}
	
	var req struct {
		FullName  *string `json:"full_name,omitempty"`
		IsActive  *bool   `json:"is_active,omitempty"`
		AvatarURL *string `json:"avatar_url,omitempty"`
		Specialty *string `json:"specialty,omitempty"`
		Bio       *string `json:"bio,omitempty"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	
	// Initialize metadata if nil
	if user.Metadata == nil {
		user.Metadata = make(map[string]interface{})
	}
	
	if req.Specialty != nil {
		user.Metadata["specialty"] = *req.Specialty
	}
	if req.Bio != nil {
		user.Metadata["bio"] = *req.Bio
	}
	
	err = instructorRepo.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update instructor"})
	}
	
	return c.JSON(http.StatusOK, userToInstructor(user, 0, 0, 0.0))
}

// DeleteInstructor deletes an instructor (changes role or deactivates)
func DeleteInstructor(c echo.Context) error {
	initInstructorRepos()
	
	id := c.Param("id")
	
	user, err := instructorRepo.GetByID(id)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Instructor not found"})
	}
	
	if user.Role != "instructor" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User is not an instructor"})
	}
	
	// Deactivate instead of delete to preserve data integrity
	user.IsActive = false
	err = instructorRepo.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete instructor"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Instructor deactivated successfully"})
}
