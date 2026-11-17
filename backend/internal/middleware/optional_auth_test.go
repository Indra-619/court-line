package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// optionalAuthRouter exposes whether userID was set by the middleware.
func optionalAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(OptionalAuthMiddleware())
	r.GET("/opt", func(c *gin.Context) {
		_, ok := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"authenticated": ok})
	})
	return r
}

func TestOptionalAuthRevokedJTITreatedAnonymous(t *testing.T) {
	defer restoreBlacklist()
	tokenBlacklist = func(ctx context.Context, jti string) (bool, error) {
		return jti == "revoked-jti", nil
	}
	r := optionalAuthRouter()

	signed := signTestJWT(t, primitive.NewObjectID(), "revoked-jti")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/opt", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("optional endpoint must stay public, got %d", w.Code)
	}
	if body := w.Body.String(); body != `{"authenticated":false}` {
		t.Errorf("revoked token must degrade to anonymous, got %s", body)
	}
}

func TestOptionalAuthLiveJTISetsUser(t *testing.T) {
	defer restoreBlacklist()
	tokenBlacklist = func(ctx context.Context, jti string) (bool, error) {
		return false, nil
	}
	r := optionalAuthRouter()

	signed := signTestJWT(t, primitive.NewObjectID(), "live-jti")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/opt", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	r.ServeHTTP(w, req)

	if body := w.Body.String(); body != `{"authenticated":true}` {
		t.Errorf("live token must authenticate, got %s", body)
	}
}

func TestOptionalAuthBlacklistErrorDegradesAnonymous(t *testing.T) {
	defer restoreBlacklist()
	tokenBlacklist = func(ctx context.Context, jti string) (bool, error) {
		return false, errors.New("mongo down")
	}
	r := optionalAuthRouter()

	signed := signTestJWT(t, primitive.NewObjectID(), "any-jti")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/opt", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("public endpoint must not 500 on blacklist outage, got %d", w.Code)
	}
	if body := w.Body.String(); body != `{"authenticated":false}` {
		t.Errorf("blacklist outage must degrade to anonymous, got %s", body)
	}
}
