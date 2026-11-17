package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Indra-619/court-line/backend/internal/handlers"
	"github.com/Indra-619/court-line/backend/internal/testutil"
)

// The lifecycle routes must exist and be auth-protected: anonymous
// requests get 401, never 404 (which would mean the route is missing).
func TestBookingLifecycleRoutesRegisteredAndProtected(t *testing.T) {
	r := SetupRouter(Deps{
		Courts:   handlers.NewCourtHandler(testutil.NewFakeCourtRepository()),
		Bookings: handlers.NewBookingHandler(testutil.NewFakeBookingRepository(), testutil.NewFakeCourtRepository()),
	})

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodDelete, "/api/bookings/64b7f1a2c3d4e5f607182930"},
		{http.MethodGet, "/api/admin/bookings"},
		{http.MethodPatch, "/api/admin/bookings/64b7f1a2c3d4e5f607182930/status"},
	}
	for _, tc := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(tc.method, tc.path, nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: expected 401 (auth-protected), got %d", tc.method, tc.path, w.Code)
		}
	}
}
