package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	// Requests per window
	Rate int
	// Window duration
	Window time.Duration
	// Key function to identify clients
	KeyFunc func(c echo.Context) string
}

// RateLimiter stores rate limit data per client
type RateLimiter struct {
	config   RateLimiterConfig
	clients  map[string]*clientData
	mu       sync.RWMutex
	stopChan chan struct{}
}

type clientData struct {
	count     int
	resetTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		config:   config,
		clients:  make(map[string]*clientData),
		stopChan: make(chan struct{}),
	}

	// Cleanup old entries periodically
	go rl.cleanup()

	return rl
}

// cleanup removes expired entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, data := range rl.clients {
				if now.After(data.resetTime) {
					delete(rl.clients, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopChan:
			return
		}
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := rl.config.KeyFunc(c)

			rl.mu.Lock()
			client, exists := rl.clients[key]
			now := time.Now()

			if !exists || now.After(client.resetTime) {
				// New client or window expired
				rl.clients[key] = &clientData{
					count:     1,
					resetTime: now.Add(rl.config.Window),
				}
				rl.mu.Unlock()
				return next(c)
			}

			// Check if rate limit exceeded
			if client.count >= rl.config.Rate {
				rl.mu.Unlock()
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"error":       "Too many requests",
					"retry_after": int(client.resetTime.Sub(now).Seconds()),
				})
			}

			// Increment count
			client.count++
			rl.mu.Unlock()

			return next(c)
		}
	}
}

// DefaultKeyFunc uses client IP as key
func DefaultKeyFunc(c echo.Context) string {
	return c.RealIP()
}

// ===== Pre-configured Rate Limiters =====

// CheckoutRateLimiter - Strict: 5 requests per minute per IP
var CheckoutRateLimiter = NewRateLimiter(RateLimiterConfig{
	Rate:    5,
	Window:  1 * time.Minute,
	KeyFunc: DefaultKeyFunc,
})

// TrackingRateLimiter - Moderate: 30 requests per minute per IP
var TrackingRateLimiter = NewRateLimiter(RateLimiterConfig{
	Rate:    30,
	Window:  1 * time.Minute,
	KeyFunc: DefaultKeyFunc,
})

// PublicAPIRateLimiter - Relaxed: 60 requests per minute per IP
var PublicAPIRateLimiter = NewRateLimiter(RateLimiterConfig{
	Rate:    60,
	Window:  1 * time.Minute,
	KeyFunc: DefaultKeyFunc,
})
