package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Simple in-memory rate limiter
type RateLimiter struct {
	requests map[string]*clientInfo
	max      int
	window   time.Duration
	mu       sync.RWMutex
}

type clientInfo struct {
	count     int
	resetTime time.Time
	mu        sync.Mutex
}

func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientInfo),
		max:      max,
		window:   window,
	}

	// Cleanup goroutine
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client IP
		clientIP := c.IP()

		rl.mu.RLock()
		info, exists := rl.requests[clientIP]
		rl.mu.RUnlock()

		if !exists {
			info = &clientInfo{
				count:     0,
				resetTime: time.Now().Add(rl.window),
			}
			rl.mu.Lock()
			rl.requests[clientIP] = info
			rl.mu.Unlock()
		}

		info.mu.Lock()
		defer info.mu.Unlock()

		// Reset if window expired
		if time.Now().After(info.resetTime) {
			info.count = 0
			info.resetTime = time.Now().Add(rl.window)
		}

		// Check rate limit
		if info.count >= rl.max {
			c.Set("X-RateLimit-Limit", strconv.Itoa(rl.max))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", info.resetTime.Format(time.RFC3339))

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
		}

		// Increment counter
		info.count++

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(rl.max))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(rl.max-info.count))
		c.Set("X-RateLimit-Reset", info.resetTime.Format(time.RFC3339))

		return c.Next()
	}
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, info := range rl.requests {
			info.mu.Lock()
			if time.Now().After(info.resetTime.Add(rl.window)) {
				delete(rl.requests, ip)
			}
			info.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}
