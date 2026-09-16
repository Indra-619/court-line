package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Indra-619/court-line/backend/internal/handlers"
	"github.com/Indra-619/court-line/backend/internal/testutil"
)

func TestHealthRoute(t *testing.T) {
	// Setup router with an in-memory court repository
	r := SetupRouter(Deps{
		Courts: handlers.NewCourtHandler(testutil.NewFakeCourtRepository()),
	})

	// Create a response recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Assert body
	expected := `{"status":"OK"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}
