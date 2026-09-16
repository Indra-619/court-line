package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
)

// exchangeForRefresh runs the exchange flow and returns the refresh
// token issued alongside the JWT.
func exchangeForRefresh(t *testing.T, h *AuthHandler, userID string) string {
	t.Helper()
	r := authTestRouter(h, primitive.NewObjectID())

	code := createExchangeCode(userID)
	body := `{"code":"` + code + `"}`
	req := httptest.NewRequest(http.MethodPost, "/exchange", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("exchange expected 200, got %d (body %s)", w.Code, w.Body.String())
	}
	var parsed struct {
		Data struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refreshToken"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid exchange JSON: %v", err)
	}
	if parsed.Data.Token == "" {
		t.Fatal("expected data.token in exchange response")
	}
	if parsed.Data.RefreshToken == "" {
		t.Fatal("expected data.refreshToken in exchange response")
	}
	return parsed.Data.RefreshToken
}

func postRefresh(t *testing.T, h *AuthHandler, refreshToken string) *httptest.ResponseRecorder {
	t.Helper()
	r := authTestRouter(h, primitive.NewObjectID())
	body := `{"refreshToken":"` + refreshToken + `"}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestExchangeIssuesRefreshTokenStoredAsHash(t *testing.T) {
	h, _, refreshTokens, _ := newSessionTestHandler()
	userID := primitive.NewObjectID().Hex()

	raw := exchangeForRefresh(t, h, userID)

	if len(raw) != 64 {
		t.Errorf("expected 64-char refresh token, got %d", len(raw))
	}
	if len(refreshTokens.Tokens) != 1 {
		t.Fatalf("expected 1 stored refresh token, got %d", len(refreshTokens.Tokens))
	}
	hash := sha256Hex(raw)
	record, ok := refreshTokens.Tokens[hash]
	if !ok {
		t.Fatal("expected the SHA-256 hash of the raw token as storage key")
	}
	if record.UserID != userID {
		t.Errorf("expected stored UserID %s, got %s", userID, record.UserID)
	}
	if strings.Contains(record.TokenHash, raw) {
		t.Error("raw token must never be persisted")
	}
	if record.ExpiresAt.Sub(record.CreatedAt) != 30*24*time.Hour {
		t.Errorf("expected 30-day TTL, got %v", record.ExpiresAt.Sub(record.CreatedAt))
	}
}

func TestRefreshHappyPathRotatesToken(t *testing.T) {
	h, users, _, _ := newSessionTestHandler()
	userID := primitive.NewObjectID()
	users.Users[userID.Hex()] = &entity.User{
		ID:    userID.Hex(),
		Email: "player@example.com",
		Name:  "Player One",
		Role:  "user",
	}

	old := exchangeForRefresh(t, h, userID.Hex())
	w := postRefresh(t, h, old)

	if w.Code != http.StatusOK {
		t.Fatalf("refresh expected 200, got %d (body %s)", w.Code, w.Body.String())
	}
	var parsed struct {
		Data struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refreshToken"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid refresh JSON: %v", err)
	}
	if parsed.Data.Token == "" || parsed.Data.RefreshToken == "" {
		t.Fatal("expected both token and refreshToken in refresh response")
	}
	if parsed.Data.RefreshToken == old {
		t.Error("refresh token must rotate to a new value")
	}

	// Replay of the old token must fail: rotation deletes it.
	if w := postRefresh(t, h, old); w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 on replay of rotated token, got %d", w.Code)
	}

	// The new token works.
	if w := postRefresh(t, h, parsed.Data.RefreshToken); w.Code != http.StatusOK {
		t.Errorf("expected 200 with the rotated token, got %d", w.Code)
	}
}

func TestRefreshUnknownTokenRejected(t *testing.T) {
	h, _, _, _ := newSessionTestHandler()
	w := postRefresh(t, h, strings.Repeat("a", 64))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Invalid or expired refresh token") {
		t.Errorf("unexpected body %s", w.Body.String())
	}
}

func TestRefreshExpiredTokenRejected(t *testing.T) {
	h, _, refreshTokens, _ := newSessionTestHandler()
	raw, err := generateRefreshToken()
	if err != nil {
		t.Fatalf("generateRefreshToken failed: %v", err)
	}
	refreshTokens.Tokens[sha256Hex(raw)] = &entity.RefreshToken{
		ID:        "expired-1",
		UserID:    primitive.NewObjectID().Hex(),
		TokenHash: sha256Hex(raw),
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-31 * 24 * time.Hour),
	}

	if w := postRefresh(t, h, raw); w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", w.Code)
	}
}

func TestRefreshRevokedTokenRejected(t *testing.T) {
	h, _, refreshTokens, _ := newSessionTestHandler()
	raw, err := generateRefreshToken()
	if err != nil {
		t.Fatalf("generateRefreshToken failed: %v", err)
	}
	refreshTokens.Tokens[sha256Hex(raw)] = &entity.RefreshToken{
		ID:        "revoked-1",
		UserID:    primitive.NewObjectID().Hex(),
		TokenHash: sha256Hex(raw),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		Revoked:   true,
	}

	if w := postRefresh(t, h, raw); w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for revoked token, got %d", w.Code)
	}
}

func TestRefreshMissingUserRejected(t *testing.T) {
	h, _, _, _ := newSessionTestHandler()
	// Exchange issues a token for a user that no longer exists.
	raw := exchangeForRefresh(t, h, primitive.NewObjectID().Hex())
	if w := postRefresh(t, h, raw); w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when user is gone, got %d", w.Code)
	}
}

func TestRefreshInvalidBodyRejected(t *testing.T) {
	h, _, _, _ := newSessionTestHandler()
	r := authTestRouter(h, primitive.NewObjectID())
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing refreshToken, got %d", w.Code)
	}
}

func TestLogoutRevokesRefreshTokensAndBlacklistsJTI(t *testing.T) {
	h, users, refreshTokens, revokedTokens := newSessionTestHandler()
	userID := primitive.NewObjectID()
	users.Users[userID.Hex()] = &entity.User{
		ID:    userID.Hex(),
		Email: "player@example.com",
		Name:  "Player One",
		Role:  "user",
	}

	raw := exchangeForRefresh(t, h, userID.Hex())
	if len(refreshTokens.Tokens) != 1 {
		t.Fatalf("expected a stored refresh token before logout")
	}

	// Simulate AuthMiddleware context: userID, jti, tokenExp.
	jti := "jti-under-test"
	exp := time.Now().Add(24 * time.Hour)
	r := newLogoutRouter(h, userID, jti, exp)
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("logout expected 200, got %d (body %s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Logged out successfully") {
		t.Errorf("unexpected body %s", w.Body.String())
	}
	if len(refreshTokens.Tokens) != 0 {
		t.Errorf("expected all refresh tokens deleted, got %d left", len(refreshTokens.Tokens))
	}
	if _, ok := revokedTokens.JTIs[jti]; !ok {
		t.Error("expected the JWT jti to be blacklisted")
	}

	// The revoked refresh token can no longer be used.
	if w := postRefresh(t, h, raw); w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 refreshing with a logged-out token, got %d", w.Code)
	}
}

func TestLogoutWithoutUserIDUnauthorized(t *testing.T) {
	h, _, _, _ := newSessionTestHandler()
	// Router with no userID in context: logout must refuse.
	r := newBareLogoutRouter(h)
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without userID, got %d", w.Code)
	}
}

// newLogoutRouter builds a router whose middleware mimics an
// authenticated request: userID, jti and tokenExp in the context.
func newLogoutRouter(h *AuthHandler, userID primitive.ObjectID, jti string, exp time.Time) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		if jti != "" {
			c.Set("jti", jti)
			c.Set("tokenExp", exp)
		}
		c.Next()
	})
	r.POST("/logout", h.Logout)
	return r
}

// newBareLogoutRouter builds a logout route with no auth context at all.
func newBareLogoutRouter(h *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/logout", h.Logout)
	return r
}
