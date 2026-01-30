package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthRoute(t *testing.T) {
	// Setup router
	r := SetupRouter()

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
