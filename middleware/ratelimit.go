package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"contrihub/database"

	"github.com/gin-gonic/gin"
)

// ─── In-memory fallback (used when Redis is unavailable) ─────────────────────

type tokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

func (tb *tokenBucket) allow() bool {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefill = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// ─── RateLimiter ─────────────────────────────────────────────────────────────

// RateLimiter provides per-token (GitHub access token) rate limiting for AI endpoints.
// Uses Redis when available, falls back to in-memory.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
	limit   int
	window  time.Duration
}

// NewRateLimiter creates a rate limiter with the given limit per window.
func NewRateLimiter(limit int, windowSeconds int) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*tokenBucket),
		limit:   limit,
		window:  time.Duration(windowSeconds) * time.Second,
	}

	// Cleanup stale in-memory entries every 5 minutes (fallback only)
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	cutoff := time.Now().Add(-10 * time.Minute)
	for key, bucket := range rl.buckets {
		if bucket.lastRefill.Before(cutoff) {
			delete(rl.buckets, key)
		}
	}
}

// getKey extracts the identifier for rate limiting.
// Prefers the Authorization token; falls back to client IP.
func getKey(c *gin.Context) string {
	if token := c.GetHeader("Authorization"); token != "" {
		return "token:" + token
	}
	return "ip:" + c.ClientIP()
}

// allowRedis checks rate limit using Redis INCR + EXPIRE (sliding window counter).
func (rl *RateLimiter) allowRedis(key string) (bool, error) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	// Increment the counter
	count, err := database.RedisClient.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	// If this is the first request in the window, set TTL
	if count == 1 {
		database.RedisClient.Expire(ctx, redisKey, rl.window)
	}

	return count <= int64(rl.limit), nil
}

// allowInMemory checks rate limit using in-memory token bucket (fallback).
func (rl *RateLimiter) allowInMemory(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	if !exists {
		bucket = &tokenBucket{
			tokens:     float64(rl.limit),
			maxTokens:  float64(rl.limit),
			refillRate: float64(rl.limit) / rl.window.Seconds(),
			lastRefill: time.Now(),
		}
		rl.buckets[key] = bucket
	}
	return bucket.allow()
}

// Middleware returns a Gin middleware that enforces rate limits.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getKey(c)
		var allowed bool

		// Try Redis first, fall back to in-memory
		if database.RedisClient != nil {
			var err error
			allowed, err = rl.allowRedis(key)
			if err != nil {
				// Redis error — gracefully fall back to in-memory
				allowed = rl.allowInMemory(key)
			}
		} else {
			allowed = rl.allowInMemory(key)
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many AI requests. Please wait before trying again.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
