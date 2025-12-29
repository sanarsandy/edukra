package middleware

import (
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
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			
			// Skip JWT for public campaign routes
			if strings.HasPrefix(path, "/api/c/") {
				return true
			}
			
			// Skip JWT for public webinar routes  
			if strings.HasPrefix(path, "/api/webinars/") ||
			   strings.HasPrefix(path, "/api/courses/") && strings.HasSuffix(path, "/webinars") {
				return true
			}
			
			// Skip JWT for public routes
			publicPaths := []string{
				"/api/courses",
				"/api/categories",
				"/api/campaign-checkout",
				"/api/transaction-status/",
				"/api/public/",
				"/api/settings",
				"/api/auth/",
				"/api/webhooks/",
				"/health",
			}
			
			for _, p := range publicPaths {
				if strings.HasPrefix(path, p) {
					return true
				}
			}
			
			return false
		},
	}
	return echojwt.WithConfig(config)
}
