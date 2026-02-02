package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/services"
	customMiddleware "github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// Validation helpers
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// isValidPassword checks password strength
// Requirements: min 8 chars, 1 uppercase, 1 lowercase, 1 digit, 1 special char
func isValidPassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password minimal 8 karakter"
	}
	
	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>_\-+=\[\]\\;'/~]`).MatchString(password)
	)
	
	if !hasUpper {
		return false, "Password harus mengandung huruf besar (A-Z)"
	}
	if !hasLower {
		return false, "Password harus mengandung huruf kecil (a-z)"
	}
	if !hasNumber {
		return false, "Password harus mengandung angka (0-9)"
	}
	if !hasSpecial {
		return false, "Password harus mengandung karakter spesial (!@#$%^&*)"
	}
	
	return true, ""
}

// Register handles user registration
func Register(c echo.Context) error {
	req := new(models.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	// Validate email
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email wajib diisi"})
	}
	if !isValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format email tidak valid"})
	}

	// Validate password strength
	if valid, errMsg := isValidPassword(req.Password); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	// Validate full name
	req.FullName = strings.TrimSpace(req.FullName)
	if req.FullName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nama lengkap wajib diisi"})
	}

	// Check if email already exists
	var existingID string
	checkQuery := `SELECT id FROM users WHERE email = $1 LIMIT 1`
	err := db.DB.QueryRow(checkQuery, req.Email).Scan(&existingID)
	if err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email sudah terdaftar"})
	} else if err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memeriksa email"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memproses password"})
	}

	// Insert user with role 'student'
	query := `INSERT INTO users (email, password_hash, full_name, role, auth_provider) 
	          VALUES ($1, $2, $3, 'student', 'email') 
	          RETURNING id, created_at`
	
	var user models.User
	user.Email = req.Email
	user.FullName = req.FullName
	user.Role = "student"
	user.AuthProvider = "email"

	err = db.DB.QueryRow(query, req.Email, string(hashedPassword), req.FullName).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mendaftarkan user: " + err.Error()})
	}

	// Generate JWT token for immediate login
	token, err := generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		// User created but token failed - still return success
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Registrasi berhasil, silakan login",
			"user":    user,
		})
	}

	return c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login handles user authentication
func Login(c echo.Context) error {
	req := new(models.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	// Validate input
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email dan password wajib diisi"})
	}

	// Find user
	var user models.User
	var googleID sql.NullString
	var authProvider sql.NullString
	
	query := `SELECT id, email, password_hash, full_name, role, google_id, auth_provider 
	          FROM users 
	          WHERE email = $1 AND (auth_provider = 'email' OR auth_provider = 'both')`
	
	err := db.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, 
		&user.Role, &googleID, &authProvider,
	)
	
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Email atau password salah"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Terjadi kesalahan sistem"})
	}

	// Set optional fields
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "email"
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Email atau password salah"})
	}

	// Generate JWT
	token, err := generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat token"})
	}

	return c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// AdminLogin handles admin-specific authentication
func AdminLogin(c echo.Context) error {
	req := new(models.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	// Validate input
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email dan password wajib diisi"})
	}

	// Find admin user
	var user models.User
	
	query := `SELECT id, email, password_hash, full_name, role 
	          FROM users 
	          WHERE email = $1 AND role = 'admin' AND (auth_provider = 'email' OR auth_provider = 'both')`
	
	err := db.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role,
	)
	
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Akses ditolak"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Terjadi kesalahan sistem"})
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Akses ditolak"})
	}

	user.AuthProvider = "email"

	// Generate JWT with admin permissions
	token, err := generateAdminToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat token"})
	}

	return c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// InstructorLogin handles instructor-specific authentication
func InstructorLogin(c echo.Context) error {
	req := new(models.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	// Validate input
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email dan password wajib diisi"})
	}

	// Find instructor user
	var user models.User
	var isActive bool
	
	query := `SELECT id, email, password_hash, full_name, role, COALESCE(is_active, true)
	          FROM users 
	          WHERE email = $1 AND role = 'instructor' AND (auth_provider = 'email' OR auth_provider = 'both')`
	
	err := db.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &isActive,
	)
	
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Email atau password salah"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Terjadi kesalahan sistem"})
	}

	// Check if instructor is active
	if !isActive {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Akun Anda belum aktif. Silakan hubungi admin."})
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Email atau password salah"})
	}

	user.AuthProvider = "email"

	// Generate JWT with instructor permissions
	token, err := generateInstructorToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat token"})
	}

	return c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// generateInstructorToken creates a JWT token for instructor users
func generateInstructorToken(userID, email, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	claims := jwt.MapClaims{
		"user_id":       userID,
		"email":         email,
		"role":          role,
		"is_instructor": true,
		"permissions":   customMiddleware.GetPermissionsForRole(role),
		"exp":           time.Now().Add(time.Hour * 4).Unix(), // 4 hours
		"iat":           time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// generateToken creates a JWT token for regular users
func generateToken(userID, email, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development only
	}

	claims := jwt.MapClaims{
		"user_id":     userID,
		"email":       email,
		"role":        role,
		"permissions": customMiddleware.GetPermissionsForRole(role),
		"exp":         time.Now().Add(time.Hour * 4).Unix(), // 4 hours
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// generateAdminToken creates a JWT token for admin users with longer expiry
func generateAdminToken(userID, email, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development only
	}

	claims := jwt.MapClaims{
		"user_id":     userID,
		"email":       email,
		"role":        role,
		"is_admin":    true,
		"permissions": customMiddleware.GetPermissionsForRole(role),
		"exp":         time.Now().Add(time.Hour * 4).Unix(), // 4 hours
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// RefreshToken refreshes an existing token
func RefreshToken(c echo.Context) error {
	userID, role, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Get email from token
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	email, _ := claims["email"].(string)

	// Generate new token
	newToken, err := generateToken(userID, email, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": newToken,
	})
}

// GetMe returns current user profile
func GetMe(c echo.Context) error {
	userID, _, err := customMiddleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var user models.User
	var googleID sql.NullString
	var authProvider sql.NullString
	var bio sql.NullString
	var phone sql.NullString
	
	query := `SELECT id, email, full_name, role, google_id, auth_provider, bio, phone, created_at 
	          FROM users WHERE id = $1`
	
	err = db.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Role, 
		&googleID, &authProvider, &bio, &phone, &user.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User tidak ditemukan"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memuat profil"})
	}

	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	}
	if bio.Valid {
		user.Bio = &bio.String
	}
	if phone.Valid {
		user.Phone = &phone.String
	}

	return c.JSON(http.StatusOK, user)
}

// Logout handles logout (client should delete token)
func Logout(c echo.Context) error {
	// JWT is stateless, so we just return success
	// Client should delete the token from storage
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Berhasil logout",
	})
}

// generateSecureToken creates a random token for password reset
// generateSecureToken creates a random token for password reset
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// RequestPasswordReset handles password reset requests
func RequestPasswordReset(c echo.Context) error {
	type Request struct {
		Email string `json:"email"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email required"})
	}

	// Always return success to prevent email enumeration
	genericResponse := map[string]string{
		"message": "If this email is registered, you will receive password reset instructions.",
	}

	// Check if user exists
	var userID string
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err == sql.ErrNoRows {
		// User not found, return generic success
		return c.JSON(http.StatusOK, genericResponse)
	} else if err != nil {
		fmt.Printf("Error finding user %s: %v\n", email, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "System error"})
	}

	// Generate secure token
	token, err := generateSecureToken(32) // 64 hex chars
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	// Set expiration (15 minutes)
	expiry := time.Now().Add(15 * time.Minute)

	// Save token to DB
	_, err = db.DB.Exec(`
		UPDATE users 
		SET reset_password_token = $1, reset_password_expires_at = $2 
		WHERE id = $3
	`, token, expiry, userID)
	
	if err != nil {
		fmt.Printf("Error saving token for %s: %v\n", email, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Send email asynchronously
	go func() {
		err := services.Email.SendResetPasswordEmail(email, token)
		if err != nil {
			fmt.Printf("Failed to send reset email to %s: %v\n", email, err)
		}
	}()

	return c.JSON(http.StatusOK, genericResponse)
}

// ResetPassword handles password reset with token
func ResetPassword(c echo.Context) error {
	type Request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Token == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token and password required"})
	}

	// Validate password strength
	if valid, errMsg := isValidPassword(req.Password); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	// Verify token and get current auth provider
	var userID string
	var expiresAt time.Time
	var authProvider string
	
	err := db.DB.QueryRow(`
		SELECT id, reset_password_expires_at,  COALESCE(auth_provider, 'email')
		FROM users 
		WHERE reset_password_token = $1
	`, req.Token).Scan(&userID, &expiresAt, &authProvider)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or expired token"})
	} else if err != nil {
		fmt.Printf("Error verifying token: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "System error"})
	}

	if time.Now().After(expiresAt) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token expired"})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process password"})
	}

	// Determine new auth provider status
	// If currently 'google', switch to 'both' so they can login with password
	newAuthProvider := authProvider
	if authProvider == "google" {
		newAuthProvider = "both"
	}

	// Update password, clear token, and update auth_provider
	_, err = db.DB.Exec(`
		UPDATE users 
		SET password_hash = $1, 
			reset_password_token = NULL, 
			reset_password_expires_at = NULL,
			auth_provider = $2
		WHERE id = $3
	`, string(hashedPassword), newAuthProvider, userID)

	if err != nil {
		fmt.Printf("Error updating password for user %s: %v\n", userID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password successfully reset. Please login with your new password.",
	})
}
