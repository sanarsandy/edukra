package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/db"
	customMiddleware "github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/middleware"
	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// Validation helpers
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidPassword(password string) bool {
	return len(password) >= 6
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

	// Validate password
	if !isValidPassword(req.Password) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password minimal 6 karakter"})
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
		"exp":           time.Now().Add(time.Hour * 12).Unix(), // 12 hours for instructor
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
		"exp":         time.Now().Add(time.Hour * 24).Unix(), // 24 hours
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
		"exp":         time.Now().Add(time.Hour * 8).Unix(), // 8 hours for admin
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
