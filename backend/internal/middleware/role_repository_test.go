package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/testutil"
)

// TestSetUserRoleLookupWiresRepository verifies the production wiring
// path: a repository-backed lookup allows admins and denies everyone
// else with the same status codes as before.
func TestSetUserRoleLookupWiresRepository(t *testing.T) {
	defer restoreLookup()

	adminID := primitive.NewObjectID()
	regularID := primitive.NewObjectID()

	users := testutil.NewFakeUserRepository()
	users.Users[adminID.Hex()] = &entity.User{ID: adminID.Hex(), Role: "admin", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	users.Users[regularID.Hex()] = &entity.User{ID: regularID.Hex(), Role: "user", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	SetUserRoleLookup(users)

	buildRouter := func(uid primitive.ObjectID) *gin.Engine {
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("userID", uid)
			c.Next()
		})
		r.Use(AdminMiddleware())
		r.POST("/admin-only", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})
		return r
	}

	t.Run("admin allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin-only", nil)
		w := httptest.NewRecorder()
		buildRouter(adminID).ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200 for admin, got %d", w.Code)
		}
	})

	t.Run("regular user forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin-only", nil)
		w := httptest.NewRecorder()
		buildRouter(regularID).ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403 for non-admin, got %d", w.Code)
		}
	})

	t.Run("unknown user forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin-only", nil)
		w := httptest.NewRecorder()
		buildRouter(primitive.NewObjectID()).ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403 for unknown user, got %d", w.Code)
		}
	})
}
