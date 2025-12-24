package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Role constants
const (
	RoleAdmin      = "admin"
	RoleInstructor = "instructor"
	RoleStudent    = "student"
)

// Permission constants
const (
	PermUsersRead        = "users:read"
	PermUsersWrite       = "users:write"
	PermUsersDelete      = "users:delete"
	PermCoursesRead      = "courses:read"
	PermCoursesWrite     = "courses:write"
	PermCoursesDelete    = "courses:delete"
	PermLessonsRead      = "lessons:read"
	PermLessonsWrite     = "lessons:write"
	PermEnrollmentsRead  = "enrollments:read"
	PermEnrollmentsWrite = "enrollments:write"
	PermTransactionsRead = "transactions:read"
	PermTransactionsWrite = "transactions:write"
	PermSettingsRead     = "settings:read"
	PermSettingsWrite    = "settings:write"
	PermProfileRead      = "profile:read"
	PermProfileWrite     = "profile:write"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[string][]string{
	RoleAdmin: {
		PermUsersRead, PermUsersWrite, PermUsersDelete,
		PermCoursesRead, PermCoursesWrite, PermCoursesDelete,
		PermLessonsRead, PermLessonsWrite,
		PermEnrollmentsRead, PermEnrollmentsWrite,
		PermTransactionsRead, PermTransactionsWrite,
		PermSettingsRead, PermSettingsWrite,
		PermProfileRead, PermProfileWrite,
	},
	RoleInstructor: {
		PermCoursesRead, PermCoursesWrite,
		PermLessonsRead, PermLessonsWrite,
		PermEnrollmentsRead,
		PermProfileRead, PermProfileWrite,
	},
	RoleStudent: {
		PermCoursesRead,
		PermLessonsRead,
		PermEnrollmentsRead, PermEnrollmentsWrite,
		PermProfileRead, PermProfileWrite,
	},
}

// GetUserFromContext extracts user info from JWT token in context
func GetUserFromContext(c echo.Context) (userID string, role string, err error) {
	user := c.Get("user")
	if user == nil {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		// Try as jwt.MapClaims directly
		if mc, ok := token.Claims.(jwt.MapClaims); ok {
			claims = &mc
		} else {
			return "", "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid claims")
		}
	}

	userID, _ = (*claims)["user_id"].(string)
	role, _ = (*claims)["role"].(string)

	if userID == "" {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in token")
	}

	return userID, role, nil
}

// RequireRole returns middleware that requires one of the specified roles
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, userRole, err := GetUserFromContext(c)
			if err != nil {
				return err
			}

			// Check if user has required role
			for _, role := range roles {
				if strings.EqualFold(userRole, role) {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Access denied: insufficient permissions")
		}
	}
}

// RequireAdmin is a shortcut for RequireRole(RoleAdmin)
func RequireAdmin() echo.MiddlewareFunc {
	return RequireRole(RoleAdmin)
}

// RequireInstructor allows only instructor role (instructor-specific routes)
func RequireInstructor() echo.MiddlewareFunc {
	return RequireRole(RoleInstructor)
}

// RequireInstructorOrAdmin allows instructor and admin
func RequireInstructorOrAdmin() echo.MiddlewareFunc {
	return RequireRole(RoleAdmin, RoleInstructor)
}

// RequirePermission returns middleware that requires a specific permission
func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, userRole, err := GetUserFromContext(c)
			if err != nil {
				return err
			}

			// Get permissions for user's role
			permissions, exists := RolePermissions[userRole]
			if !exists {
				return echo.NewHTTPError(http.StatusForbidden, "Unknown role")
			}

			// Check if user has required permission
			for _, perm := range permissions {
				if perm == permission {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Access denied: missing permission "+permission)
		}
	}
}

// HasPermission checks if a role has a specific permission
func HasPermission(role, permission string) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, perm := range permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// GetPermissionsForRole returns all permissions for a given role
func GetPermissionsForRole(role string) []string {
	permissions, exists := RolePermissions[role]
	if !exists {
		return []string{}
	}
	return permissions
}
