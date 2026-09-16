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

// originalLookup holds the default MongoDB implementation so tests can
// restore it when finished.
var originalLookup = userRoleLookup

func setupAdminTestRouter(role string, lookupErr error, setUserID bool) *gin.Engine {
	gin.SetMode(gin.TestMode)

	userRoleLookup = func(ctx context.Context, userID primitive.ObjectID) (string, error) {
		return role, lookupErr
	}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		if setUserID {
			c.Set("userID", primitive.NewObjectID())
		}
		c.Next()
	})
	r.Use(AdminMiddleware())
	r.POST("/admin-only", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	return r
}

func restoreLookup() {
	userRoleLookup = originalLookup
}

func TestAdminMiddlewareAdminPasses(t *testing.T) {
	defer restoreLookup()
	r := setupAdminTestRouter("admin", nil, true)

	req := httptest.NewRequest(http.MethodPost, "/admin-only", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for admin user, got %d", w.Code)
	}
}

func TestAdminMiddlewareNonAdminForbidden(t *testing.T) {
	defer restoreLookup()
	r := setupAdminTestRouter("user", nil, true)

	req := httptest.NewRequest(http.MethodPost, "/admin-only", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for non-admin user, got %d", w.Code)
	}
}

func TestAdminMiddlewareMissingUserForbidden(t *testing.T) {
	defer restoreLookup()
	r := setupAdminTestRouter("", errors.New("mongo: no documents in result"), false)

	req := httptest.NewRequest(http.MethodPost, "/admin-only", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 when userID missing from context, got %d", w.Code)
	}
}

func TestAdminMiddlewareLookupErrorForbidden(t *testing.T) {
	defer restoreLookup()
	r := setupAdminTestRouter("", errors.New("mongo: no documents in result"), true)

	req := httptest.NewRequest(http.MethodPost, "/admin-only", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 when role lookup fails, got %d", w.Code)
	}
}
