package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Indra-619/court-line/backend/pkg/config"
)

var originalBlacklist = tokenBlacklist

func restoreBlacklist() {
	tokenBlacklist = originalBlacklist
}

// signTestJWT builds an HS256 token with the given jti for middleware tests.
func signTestJWT(t *testing.T, userID primitive.ObjectID, jti string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"userId": userID.Hex(),
		"exp":    time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}
	if jti != "" {
		claims["jti"] = jti
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.JWTSecret()))
	if err != nil {
		t.Fatalf("failed to sign test JWT: %v", err)
	}
	return signed
}

func setupAuthTestRouter(blacklisted map[string]bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	tokenBlacklist = func(ctx context.Context, jti string) (bool, error) {
		return blacklisted[jti], nil
	}
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestAuthMiddlewareBlacklistedJTIRejected(t *testing.T) {
	defer restoreBlacklist()
	r := setupAuthTestRouter(map[string]bool{"revoked-jti": true})

	signed := signTestJWT(t, primitive.NewObjectID(), "revoked-jti")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for blacklisted jti, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Token revoked") {
		t.Errorf("expected 'Token revoked' error, got %s", w.Body.String())
	}
}

func TestAuthMiddlewareLiveJTIPasses(t *testing.T) {
	defer restoreBlacklist()
	r := setupAuthTestRouter(map[string]bool{"other-jti": true})

	signed := signTestJWT(t, primitive.NewObjectID(), "live-jti")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for live jti, got %d (body %s)", w.Code, w.Body.String())
	}
}

func TestAuthMiddlewareLegacyTokenWithoutJTIPasses(t *testing.T) {
	defer restoreBlacklist()
	// Default no-op blacklist must keep pre-jti tokens working.
	tokenBlacklist = func(ctx context.Context, jti string) (bool, error) {
		return false, nil
	}
	r := setupAuthTestRouter(nil)

	signed := signTestJWT(t, primitive.NewObjectID(), "")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for legacy token without jti, got %d (body %s)", w.Code, w.Body.String())
	}
}

func TestAuthMiddlewareInvalidTokenRejected(t *testing.T) {
	defer restoreBlacklist()
	r := setupAuthTestRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for malformed token, got %d", w.Code)
	}
}

func TestAuthMiddlewareBlacklistErrorFailsClosed(t *testing.T) {
	defer restoreBlacklist()
	gin.SetMode(gin.TestMode)
	tokenBlacklist = func(ctx context.Context, jti string) (bool, error) {
		return false, context.DeadlineExceeded
	}
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	signed := signTestJWT(t, primitive.NewObjectID(), "any-jti")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when blacklist lookup errors, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Authentication service unavailable") {
		t.Errorf("expected unavailable error body, got %s", w.Body.String())
	}
}
