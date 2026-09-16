package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/testutil"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetCourtsReturnsDataArray(t *testing.T) {
	repo := testutil.NewFakeCourtRepository()
	repo.Courts["1"] = &entity.Court{
		ID:           "64b7f1a2c3d4e5f607182930",
		Name:         "Lapangan A",
		Type:         "Futsal",
		Location:     "Jakarta",
		Description:  "Rumput sintetis",
		PricePerHour: 150000,
		ImageURL:     "https://example.com/a.jpg",
		Facilities:   []string{"parkir", "lampu"},
		IsAvailable:  true,
	}

	h := NewCourtHandler(repo)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/api/courts", h.GetCourts)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/courts", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Data []entity.Court `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(body.Data) != 1 {
		t.Fatalf("expected 1 court, got %d", len(body.Data))
	}
	court := body.Data[0]
	if court.Name != "Lapangan A" || court.PricePerHour != 150000 || !court.IsAvailable {
		t.Errorf("unexpected court payload: %+v", court)
	}
}

func TestGetCourtsEmptyReturnsEmptyArray(t *testing.T) {
	h := NewCourtHandler(testutil.NewFakeCourtRepository())
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/api/courts", h.GetCourts)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/courts", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != `{"data":[]}` {
		t.Errorf("expected empty data array, got %s", w.Body.String())
	}
}

func TestGetCourtByIDNotFound(t *testing.T) {
	h := NewCourtHandler(testutil.NewFakeCourtRepository())
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/api/courts/:id", h.GetCourtByID)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/courts/64b7f1a2c3d4e5f607182930", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if w.Body.String() != `{"error":"Court not found"}` {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestGetCourtByIDInvalidHex(t *testing.T) {
	h := NewCourtHandler(testutil.NewFakeCourtRepository())
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/api/courts/:id", h.GetCourtByID)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/courts/not-a-hex-id", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if w.Body.String() != `{"error":"Invalid court ID"}` {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}
