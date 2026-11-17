package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimiter is a fixed-window per-IP counter. Auth endpoints are the
// abuse surface (token exchange, refresh, OAuth callback); a modest
// window keeps scripted guessing impractical without a reverse proxy.
// State is in-memory: acceptable for a single-instance deployment, and
// it fails open on restart rather than locking users out.
type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string]*rateWindow
	limit  int
	window time.Duration
	nowFn  func() time.Time
}

type rateWindow struct {
	count int
	start time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		hits:   make(map[string]*rateWindow),
		limit:  limit,
		window: window,
		nowFn:  time.Now,
	}
}

// allow reports whether the key may proceed, recording the hit when it
// may. Expired windows are reset lazily on access.
func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.nowFn()
	w, ok := rl.hits[key]
	if !ok || now.Sub(w.start) >= rl.window {
		rl.hits[key] = &rateWindow{count: 1, start: now}
		// Opportunistic cleanup so the map cannot grow unbounded under
		// IP churn; cheap because the map stays small at this limit.
		if len(rl.hits) > 1024 {
			for k, v := range rl.hits {
				if now.Sub(v.start) >= rl.window {
					delete(rl.hits, k)
				}
			}
		}
		return true
	}
	if w.count >= rl.limit {
		return false
	}
	w.count++
	return true
}

// RateLimitMiddleware caps requests per client IP over a fixed window.
// Excess requests get 429 with Retry-After guidance.
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	rl := newRateLimiter(limit, window)
	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}
