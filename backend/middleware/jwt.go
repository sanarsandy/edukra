package middleware

import (
	"log"
	"os"
	"strings"
	
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development only
	}
	
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},
		SigningKey: []byte(jwtSecret),
		ErrorHandler: func(c echo.Context, err error) error {
			authHeader := c.Request().Header.Get("Authorization")
			log.Printf("[JWT] 401 Error on %s: %v | Auth Header: %s", c.Request().URL.Path, err, authHeader)
			return echo.NewHTTPError(401, "Unauthorized")
		},
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			
			// Skip JWT for public campaign routes
			if strings.HasPrefix(path, "/api/c/") {
				return true
			}
			
			// Skip JWT for public webinar routes  
			if strings.HasPrefix(path, "/api/webinars/") {
				return true
			}
			
			// Skip for courses/:id/webinars (public)
			if strings.HasPrefix(path, "/api/courses/") && strings.HasSuffix(path, "/webinars") {
				return true
			}
			
			// EXACT match for public course routes (listing and detail only)
			// /api/courses - list all courses (public)
			if path == "/api/courses" {
				return true
			}
			
			// /api/courses/:id - get single course (public)
			// Match pattern: /api/courses/{uuid} without additional path segments
			if strings.HasPrefix(path, "/api/courses/") {
				// Check if this is just course detail (no sub-path after ID)
				parts := strings.Split(path, "/")
				// /api/courses/:id has 4 parts: ["", "api", "courses", "id"]
				if len(parts) == 4 {
					return true
				}
				// Also allow /api/courses/:id/ratings and /api/courses/:id/ratings/stats (public stats)
				if len(parts) >= 5 && parts[4] == "ratings" {
					// /api/courses/:id/ratings - public list
					// /api/courses/:id/ratings/stats - public stats
					if len(parts) == 5 || (len(parts) == 6 && parts[5] == "stats") {
						return true
					}
				}
			}
			
			// Skip JWT for other public routes (exact prefix match with trailing consideration)
			publicPrefixes := []string{
				"/api/categories",
				"/api/campaign-checkout",
				"/api/transaction-status/",
				"/api/public/",
				"/api/settings",
				"/api/auth/",
				"/api/webhooks/",
				"/health",
			}
			
			for _, p := range publicPrefixes {
				if strings.HasPrefix(path, p) {
					return true
				}
			}
			
			return false
		},
	}
	return echojwt.WithConfig(config)
}
