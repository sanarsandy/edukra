package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/domain"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

var userRepo *postgres.UserRepository

func initUserRepos() {
	if userRepo == nil && db.DB != nil {
		userRepo = postgres.NewUserRepository(db.DB)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// GetCurrentUser returns the authenticated user's profile
func GetCurrentUser(c echo.Context) error {
	initUserRepos()
	
	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}
	
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	
	return c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser updates the authenticated user's profile
func UpdateCurrentUser(c echo.Context) error {
	initUserRepos()
	
	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	
	var req struct {
		FullName  *string `json:"full_name,omitempty"`
		AvatarURL *string `json:"avatar_url,omitempty"`
		Bio       *string `json:"bio,omitempty"`
		Phone     *string `json:"phone,omitempty"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	
	err = userRepo.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}
	
	return c.JSON(http.StatusOK, user)
}

// ChangePassword allows authenticated user to change their password
func ChangePassword(c echo.Context) error {
	initUserRepos()
	
	userID := getUserIDFromToken(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}
	
	user, err := userRepo.GetByID(userID)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	// Validate inputs
	if req.CurrentPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Current password is required"})
	}
	if req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "New password is required"})
	}
	if len(req.NewPassword) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "New password must be at least 8 characters"})
	}
	
	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Password saat ini salah"})
	}
	
	// Hash new password
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	
	// Update password
	err = userRepo.UpdatePassword(userID, hashedPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Password berhasil diubah"})
}

// CreateUser creates a new user (admin only)
func CreateUser(c echo.Context) error {
	initUserRepos()
	
	var req struct {
		Email    string          `json:"email"`
		Password string          `json:"password"`
		FullName string          `json:"full_name"`
		Role     domain.UserRole `json:"role"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	// Validate required fields
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password is required"})
	}
	if req.FullName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Full name is required"})
	}
	if req.Role == "" {
		req.Role = "student"
	}
	
	// Check if email already exists
	existingUser, _ := userRepo.GetByEmail("default", req.Email)
	if existingUser != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email already registered"})
	}
	
	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	
	// Create user
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         req.Role,
		AuthProvider: "email",
		IsActive:     true,
	}
	
	err = userRepo.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user: " + err.Error()})
	}
	
	return c.JSON(http.StatusCreated, user)
}

// UserWithStats includes enrollment count
type UserWithStats struct {
	*domain.User
	EnrollmentCount int `json:"enrollment_count"`
}

// ListUsers returns all users (admin only)
func ListUsers(c echo.Context) error {
	initUserRepos()
	initEnrollmentRepos()
	
	tenantID := c.QueryParam("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	
	// Optional filter by status
	statusFilter := c.QueryParam("status") // active, inactive, all
	// Optional filter by role
	roleFilter := c.QueryParam("role") // admin, instructor, student
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}
	
	users, err := userRepo.ListByTenant(tenantID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}
	
	// Apply filters and enrich with enrollment count
	var usersWithStats []*UserWithStats
	for _, user := range users {
		// Apply status filter
		if statusFilter != "" && statusFilter != "all" {
			if statusFilter == "active" && !user.IsActive {
				continue
			}
			if statusFilter == "inactive" && user.IsActive {
				continue
			}
		}
		
		// Apply role filter
		if roleFilter != "" && string(user.Role) != roleFilter {
			continue
		}
		
		// Get enrollment count for each user
		enrollCount := 0
		if enrollmentRepo != nil {
			count, _ := enrollmentRepo.CountByUser(user.ID)
			enrollCount = count
		}
		
		usersWithStats = append(usersWithStats, &UserWithStats{
			User:            user,
			EnrollmentCount: enrollCount,
		})
	}
	
	total, _ := userRepo.CountByTenant(tenantID)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"users":  usersWithStats,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUser returns a single user by ID (admin only)
func GetUser(c echo.Context) error {
	initUserRepos()
	
	id := c.Param("id")
	user, err := userRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}
	
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	
	return c.JSON(http.StatusOK, user)
}

// UpdateUser updates a user (admin only)
func UpdateUser(c echo.Context) error {
	initUserRepos()
	
	id := c.Param("id")
	
	user, err := userRepo.GetByID(id)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	
	var req struct {
		FullName  *string          `json:"full_name,omitempty"`
		Role      *domain.UserRole `json:"role,omitempty"`
		IsActive  *bool            `json:"is_active,omitempty"`
		AvatarURL *string          `json:"avatar_url,omitempty"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	
	err = userRepo.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}
	
	return c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user (admin only)
func DeleteUser(c echo.Context) error {
	initUserRepos()
	
	id := c.Param("id")
	
	user, err := userRepo.GetByID(id)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	
	err = userRepo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
