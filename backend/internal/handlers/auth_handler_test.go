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
	"github.com/Indra-619/court-line/backend/internal/testutil"
)

func authTestRouter(h *AuthHandler, userID primitive.ObjectID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.GET("/me", h.GetCurrentUser)
	r.POST("/logout", h.Logout)
	r.POST("/exchange", h.ExchangeToken)
	return r
}

func TestGetCurrentUserReturnsUserData(t *testing.T) {
	userID := primitive.NewObjectID()
	users := testutil.NewFakeUserRepository()
	users.Users[userID.Hex()] = &entity.User{
		ID:        userID.Hex(),
		GoogleID:  "google-123",
		Email:     "player@example.com",
		Name:      "Player One",
		Picture:   "https://example.com/pic.png",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	h := NewAuthHandler(users)
	r := authTestRouter(h, userID)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
	}

	var body map[string]map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	data := body["data"]
	if data == nil {
		t.Fatal("expected data envelope")
	}
	if data["email"] != "player@example.com" {
		t.Errorf("expected email player@example.com, got %v", data["email"])
	}
	if data["role"] != "user" {
		t.Errorf("expected role user, got %v", data["role"])
	}
}

func TestGetCurrentUserNotFound(t *testing.T) {
	users := testutil.NewFakeUserRepository()
	h := NewAuthHandler(users)
	r := authTestRouter(h, primitive.NewObjectID())

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "User not found") {
		t.Errorf("expected 'User not found' error, got %s", w.Body.String())
	}
}

func TestLogoutAlwaysSucceeds(t *testing.T) {
	h := NewAuthHandler(testutil.NewFakeUserRepository())
	r := authTestRouter(h, primitive.NewObjectID())

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Logged out successfully") {
		t.Errorf("unexpected body %s", w.Body.String())
	}
}

func TestExchangeTokenRoundTrip(t *testing.T) {
	h := NewAuthHandler(testutil.NewFakeUserRepository())
	r := authTestRouter(h, primitive.NewObjectID())

	code := createExchangeCode(primitive.NewObjectID().Hex())
	body := `{"code":"` + code + `"}`
	req := httptest.NewRequest(http.MethodPost, "/exchange", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"token":"`) {
		t.Errorf("expected token in data envelope, got %s", w.Body.String())
	}

	// Replay must be rejected: codes are single-use.
	req = httptest.NewRequest(http.MethodPost, "/exchange", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 on replay, got %d", w.Code)
	}
}
