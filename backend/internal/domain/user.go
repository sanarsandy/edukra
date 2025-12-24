package domain

import (
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleInstructor UserRole = "instructor"
	RoleStudent    UserRole = "student"
)

// AuthProvider represents how the user authenticated
type AuthProvider string

const (
	AuthEmail  AuthProvider = "email"
	AuthGoogle AuthProvider = "google"
	AuthBoth   AuthProvider = "both"
)

// User represents a user in the LMS system
type User struct {
	ID           string       `json:"id"`
	TenantID     *string      `json:"tenant_id,omitempty"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"-"` // Never expose in JSON
	Role         UserRole     `json:"role"`
	FullName     string       `json:"full_name"`
	AvatarURL    *string      `json:"avatar_url,omitempty"`
	Bio          *string      `json:"bio,omitempty"`
	Phone        *string      `json:"phone,omitempty"`
	GoogleID     *string      `json:"google_id,omitempty"`
	AuthProvider AuthProvider `json:"auth_provider"`
	IsActive     bool                   `json:"is_active"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	TenantID string `json:"tenant_id,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantID string `json:"tenant_id,omitempty"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(id string) (*User, error)
	GetByEmail(tenantID, email string) (*User, error)
	GetByGoogleID(googleID string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
	ListByTenant(tenantID string, limit, offset int) ([]*User, error)
}

// AuthService defines the interface for authentication logic
type AuthService interface {
	Register(req *RegisterRequest) (*User, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	ValidateToken(token string) (*User, error)
	RefreshToken(token string) (string, error)
}
