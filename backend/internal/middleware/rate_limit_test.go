package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func rateLimitRouter(limit int, window time.Duration) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/auth/exchange", RateLimitMiddleware(limit, window), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestRateLimitBlocksAfterLimit(t *testing.T) {
	r := rateLimitRouter(3, time.Minute)
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/exchange", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/exchange", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("4th request: expected 429, got %d", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Error("expected Retry-After header on 429")
	}
}

func TestRateLimitPerIPIsolation(t *testing.T) {
	r := rateLimitRouter(1, time.Minute)
	hit := func(ip string) int {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/exchange", nil)
		req.RemoteAddr = ip + ":1234"
		r.ServeHTTP(w, req)
		return w.Code
	}
	if got := hit("10.0.0.1"); got != http.StatusOK {
		t.Fatalf("first IP first hit: expected 200, got %d", got)
	}
	if got := hit("10.0.0.2"); got != http.StatusOK {
		t.Fatalf("second IP must not be limited by first, got %d", got)
	}
	if got := hit("10.0.0.1"); got != http.StatusTooManyRequests {
		t.Fatalf("first IP second hit: expected 429, got %d", got)
	}
}

func TestRateLimitWindowResets(t *testing.T) {
	rl := newRateLimiter(1, time.Minute)
	now := time.Now()
	rl.nowFn = func() time.Time { return now }

	if !rl.allow("ip") {
		t.Fatal("first hit must pass")
	}
	if rl.allow("ip") {
		t.Fatal("second hit within window must be blocked")
	}
	now = now.Add(61 * time.Second)
	if !rl.allow("ip") {
		t.Fatal("hit after window expiry must pass")
	}
}
