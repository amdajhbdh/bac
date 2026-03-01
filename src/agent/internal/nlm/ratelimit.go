package nlm

import (
	"context"
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*Bucket
	config  map[string]time.Duration
}

type Bucket struct {
	lastRequest time.Time
	cooldown    time.Duration
}

var defaultCooldowns = map[string]time.Duration{
	"query":    2 * time.Second,
	"source":   2 * time.Second,
	"audio":    5 * time.Second,
	"research": 5 * time.Second,
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*Bucket),
		config:  defaultCooldowns,
	}
}

func (r *RateLimiter) SetCooldown(operation string, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.config[operation] = duration
}

func (r *RateLimiter) Acquire(ctx context.Context, notebookID, operation string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cooldown := r.config[operation]
	if cooldown == 0 {
		cooldown = r.config["query"]
	}

	key := notebookID + ":" + operation

	bucket, exists := r.buckets[key]
	if !exists {
		r.buckets[key] = &Bucket{
			lastRequest: time.Now().Add(-cooldown),
			cooldown:    cooldown,
		}
		return nil
	}

	elapsed := time.Since(bucket.lastRequest)
	if elapsed < cooldown {
		waitTime := cooldown - elapsed
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}
	}

	bucket.lastRequest = time.Now()
	return nil
}

func (r *RateLimiter) Release(notebookID, operation string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := notebookID + ":" + operation
	if bucket, exists := r.buckets[key]; exists {
		bucket.lastRequest = time.Now()
	}
}

func (r *RateLimiter) WaitTime(notebookID, operation string) time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := notebookID + ":" + operation
	cooldown := r.config[operation]
	if cooldown == 0 {
		cooldown = r.config["query"]
	}

	bucket, exists := r.buckets[key]
	if !exists {
		return 0
	}

	elapsed := time.Since(bucket.lastRequest)
	if elapsed < cooldown {
		return cooldown - elapsed
	}
	return 0
}
