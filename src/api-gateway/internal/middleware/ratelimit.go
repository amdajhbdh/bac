package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements simple in-memory rate limiting
type RateLimiter struct {
	mu       sync.RWMutex
	clients  map[string]*clientLimiter
	window   time.Duration
	maxCount int
}

type clientLimiter struct {
	count    int
	windowStart time.Time
}

func NewRateLimiter(window time.Duration, maxCount int) *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string]*clientLimiter),
		window:   window,
		maxCount: maxCount,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[key]
	if !exists {
		rl.clients[key] = &clientLimiter{
			count: 1,
			windowStart: now,
		}
		return true
	}

	// Reset window if expired
	if now.Sub(client.windowStart) > rl.window {
		client.count = 0
		client.windowStart = now
	}

	if client.count >= rl.maxCount {
		return false
	}

	client.count++
	return true
}

// RateLimit middleware constructor
func RateLimit(window time.Duration, maxCount int) gin.HandlerFunc {
	limiter := NewRateLimiter(window, maxCount)

	return func(c *gin.Context) {
		key := c.ClientIP()
		if !limiter.Allow(key) {
			c.JSON(429, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}