package middleware

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// RateLimiter provides IP-based rate limiting using an in-memory sliding window.
type RateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	maxCount int
	window   time.Duration
}

// NewRateLimiter creates a rate limiter with maxCount requests per window.
func NewRateLimiter(maxCount int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string][]time.Time),
		maxCount: maxCount,
		window:   window,
	}
}

// Allow checks if a request from the given key (IP) is allowed.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Remove old entries
	attempts := rl.attempts[key]
	valid := make([]time.Time, 0, len(attempts))
	for _, t := range attempts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.maxCount {
		rl.attempts[key] = valid
		return false
	}

	rl.attempts[key] = append(valid, now)
	return true
}

// RateLimit returns a middleware that limits requests by client IP.
func RateLimit(maxCount int, window time.Duration) func(r *ghttp.Request) {
	limiter := NewRateLimiter(maxCount, window)
	return func(r *ghttp.Request) {
		ip := r.GetClientIp()
		if !limiter.Allow(ip) {
			g.Log().Warningf(r.GetCtx(), "rate limit exceeded for IP: %s", ip)
			r.Response.WriteJsonExit(g.Map{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			return
		}
		r.Middleware.Next()
	}
}
