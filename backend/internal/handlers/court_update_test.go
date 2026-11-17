package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/testutil"
)

func courtUpdateRouter(h *CourtHandler) *gin.Engine {
	r := gin.New()
	r.PUT("/api/courts/:id", h.UpdateCourt)
	return r
}

func putCourt(t *testing.T, r *gin.Engine, id, body string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/courts/"+id, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func TestUpdateCourtCanDisableAvailability(t *testing.T) {
	repo := testutil.NewFakeCourtRepository()
	id := "64b7f1a2c3d4e5f607182930"
	repo.Courts[id] = &entity.Court{ID: id, Name: "A", Type: "Futsal", Location: "JKT", PricePerHour: 100000, IsAvailable: true}
	r := courtUpdateRouter(NewCourtHandler(repo))

	body := `{"name":"A","type":"Futsal","location":"JKT","pricePerHour":100000,"isAvailable":false}`
	w := putCourt(t, r, id, body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}
	if repo.Courts[id].IsAvailable {
		t.Error("expected court to be disabled after explicit isAvailable=false")
	}
}

func TestUpdateCourtOmittedAvailabilityPreserved(t *testing.T) {
	repo := testutil.NewFakeCourtRepository()
	id := "64b7f1a2c3d4e5f607182930"
	repo.Courts[id] = &entity.Court{ID: id, Name: "A", Type: "Futsal", Location: "JKT", PricePerHour: 100000, IsAvailable: true}
	r := courtUpdateRouter(NewCourtHandler(repo))

	body := `{"name":"B","type":"Futsal","location":"JKT","pricePerHour":120000}`
	w := putCourt(t, r, id, body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !repo.Courts[id].IsAvailable {
		t.Error("omitted isAvailable must preserve the stored true value")
	}
	if repo.Courts[id].Name != "B" {
		t.Errorf("expected name updated to B, got %s", repo.Courts[id].Name)
	}
}

func TestUpdateCourtNotFound(t *testing.T) {
	r := courtUpdateRouter(NewCourtHandler(testutil.NewFakeCourtRepository()))
	w := putCourt(t, r, "64b7f1a2c3d4e5f607182930", `{"name":"A","type":"Futsal","location":"JKT","pricePerHour":1}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreateBookingUnavailableCourtRejected(t *testing.T) {
	courts, bookings, courtID, userID := newBookingFakes(t)
	courts.Courts[courtID].IsAvailable = false
	h := NewBookingHandler(bookings, courts)
	r := bookingTestRouter(h, userID)

	date := futureDate(t)
	body := `{"courtId":"` + courtID + `","customerName":"Budi","customerPhone":"081234567890","date":"` + date + `","startTime":"18:00","endTime":"20:00"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 for unavailable court, got %d (body: %s)", w.Code, w.Body.String())
	}
	if len(bookings.Bookings) != 0 {
		t.Error("no booking should be stored for an unavailable court")
	}
}

func TestGetCourtsPagination(t *testing.T) {
	repo := testutil.NewFakeCourtRepository()
	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("64b7f1a2c3d4e5f6071829%02d", i)
		repo.Courts[id] = &entity.Court{ID: id, Name: fmt.Sprintf("Court %d", i)}
	}
	h := NewCourtHandler(repo)

	get := func(query string) []entity.Court {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.GET("/api/courts", h.GetCourts)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/courts"+query, nil)
		r.ServeHTTP(w, c.Request)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var body struct {
			Data []entity.Court `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		return body.Data
	}

	if got := get("?limit=2&offset=0"); len(got) != 2 {
		t.Errorf("limit=2: expected 2 courts, got %d", len(got))
	}
	if got := get("?limit=2&offset=4"); len(got) != 1 {
		t.Errorf("offset=4 of 5: expected 1 court, got %d", len(got))
	}
	if got := get("?offset=99"); len(got) != 0 {
		t.Errorf("offset past end: expected empty array, got %d", len(got))
	}
	if got := get("?limit=9999"); len(got) != 5 {
		// 5 stored courts, cap 100 — all fit; cap itself exercised below.
		t.Errorf("expected all 5 courts, got %d", len(got))
	}
	if got := get("?limit=abc"); len(got) != 5 {
		t.Errorf("invalid limit must fall back to default, got %d", len(got))
	}
}

func TestListPaginationCap(t *testing.T) {
	repo := testutil.NewFakeCourtRepository()
	for i := 0; i < 150; i++ {
		id := fmt.Sprintf("64b7f1a2c3d4e5f60718%04d", i)
		repo.Courts[id] = &entity.Court{ID: id}
	}
	h := NewCourtHandler(repo)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/api/courts", h.GetCourts)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/courts?limit=500", nil)
	r.ServeHTTP(w, c.Request)

	var body struct {
		Data []entity.Court `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Data) != maxListLimit {
		t.Errorf("expected cap at %d, got %d", maxListLimit, len(body.Data))
	}
}
