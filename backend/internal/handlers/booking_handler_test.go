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

func futureDate(t *testing.T) string {
	t.Helper()
	return time.Now().Add(48 * time.Hour).Format("2006-01-02")
}

// bookingTestRouter wires a BookingHandler behind a stub auth
// middleware that sets the userID the real middleware would set.
func bookingTestRouter(h *BookingHandler, userID primitive.ObjectID) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/api/bookings", h.CreateBooking)
	r.GET("/api/bookings", h.GetBookings)
	r.GET("/api/courts/:id/bookings", h.GetBookingsByCourtID)
	return r
}

func newBookingFakes(t *testing.T) (*testutil.FakeCourtRepository, *testutil.FakeBookingRepository, string, primitive.ObjectID) {
	t.Helper()
	courtID := "64b7f1a2c3d4e5f607182930"
	courts := testutil.NewFakeCourtRepository()
	courts.Courts[courtID] = &entity.Court{
		ID:           courtID,
		Name:         "Lapangan A",
		PricePerHour: 100000,
		IsAvailable:  true,
	}
	bookings := testutil.NewFakeBookingRepository()
	userID := primitive.NewObjectID()
	return courts, bookings, courtID, userID
}

func TestCreateBookingHappyPath(t *testing.T) {
	courts, bookings, courtID, userID := newBookingFakes(t)
	h := NewBookingHandler(bookings, courts)
	r := bookingTestRouter(h, userID)

	date := futureDate(t)
	body := `{"courtId":"` + courtID + `","customerName":"Budi","customerPhone":"081234567890","date":"` + date + `","startTime":"18:00","endTime":"20:00"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body: %s)", w.Code, w.Body.String())
	}

	var resp struct {
		Data entity.Booking `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	created := resp.Data
	if created.TotalPrice != 200000 {
		t.Errorf("expected price 200000 for 2h x 100000, got %v", created.TotalPrice)
	}
	if created.Status != entity.BookingStatusPending {
		t.Errorf("expected pending status, got %s", created.Status)
	}
	if created.UserID != userID.Hex() {
		t.Errorf("expected userId %s, got %s", userID.Hex(), created.UserID)
	}
	if created.ID == "" {
		t.Error("expected generated booking ID")
	}
	if len(bookings.Bookings) != 1 {
		t.Errorf("expected 1 stored booking, got %d", len(bookings.Bookings))
	}
}

func TestCreateBookingConflict(t *testing.T) {
	courts, bookings, courtID, userID := newBookingFakes(t)
	date := futureDate(t)
	// Seed an active booking that overlaps 18:00-20:00
	bookings.Bookings = append(bookings.Bookings, &entity.Booking{
		ID:        "existing",
		CourtID:   courtID,
		UserID:    primitive.NewObjectID().Hex(),
		Date:      date,
		StartTime: "19:00",
		EndTime:   "21:00",
		Status:    entity.BookingStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	h := NewBookingHandler(bookings, courts)
	r := bookingTestRouter(h, userID)

	body := `{"courtId":"` + courtID + `","customerName":"Budi","customerPhone":"081234567890","date":"` + date + `","startTime":"18:00","endTime":"20:00"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d (body: %s)", w.Code, w.Body.String())
	}

	var resp struct {
		Error    string `json:"error"`
		Conflict struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		} `json:"conflict"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !strings.Contains(resp.Error, "Time slot conflicts") {
		t.Errorf("unexpected error message: %s", resp.Error)
	}
	if resp.Conflict.StartTime != "19:00" || resp.Conflict.EndTime != "21:00" {
		t.Errorf("unexpected conflict payload: %+v", resp.Conflict)
	}
	if len(bookings.Bookings) != 1 {
		t.Errorf("conflicting booking must not be stored, got %d bookings", len(bookings.Bookings))
	}
}

func TestCreateBookingCourtNotFound(t *testing.T) {
	courts, bookings, _, userID := newBookingFakes(t)
	h := NewBookingHandler(bookings, courts)
	r := bookingTestRouter(h, userID)

	date := futureDate(t)
	body := `{"courtId":"64b7f1a2c3d4e5f607182939","customerName":"Budi","customerPhone":"081234567890","date":"` + date + `","startTime":"18:00","endTime":"20:00"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if w.Body.String() != `{"error":"Court not found"}` {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestGetBookingsByCourtIDRedactsPII(t *testing.T) {
	courts, bookings, courtID, userID := newBookingFakes(t)
	bookings.Bookings = append(bookings.Bookings, &entity.Booking{
		ID:            "64b7f1a2c3d4e5f607182940",
		CourtID:       courtID,
		UserID:        userID.Hex(),
		CustomerName:  "Budi Santoso",
		CustomerPhone: "081234567890",
		Date:          futureDate(t),
		StartTime:     "18:00",
		EndTime:       "20:00",
		TotalPrice:    200000,
		Status:        entity.BookingStatusConfirmed,
	})

	h := NewBookingHandler(bookings, courts)
	r := bookingTestRouter(h, userID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/courts/"+courtID+"/bookings", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	out := w.Body.String()
	for _, forbidden := range []string{"Budi Santoso", "081234567890", "customerName", "customerPhone", "totalPrice"} {
		if strings.Contains(out, forbidden) {
			t.Errorf("public court schedule leaks %q: %s", forbidden, out)
		}
	}
	if !strings.Contains(out, `"status":"confirmed"`) {
		t.Errorf("expected status field in response: %s", out)
	}
}
