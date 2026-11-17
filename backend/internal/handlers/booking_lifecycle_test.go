package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/domain/repository"
	"github.com/Indra-619/court-line/backend/internal/testutil"
)

func seedLifecycleBooking(t *testing.T, status entity.BookingStatus) (*testutil.FakeBookingRepository, string, string) {
	t.Helper()
	bookings := testutil.NewFakeBookingRepository()
	courts := testutil.NewFakeCourtRepository()
	userID := "64b7f1a2c3d4e5f607182931"
	bookingID := "64b7f1a2c3d4e5f607182940"
	bookings.Bookings = append(bookings.Bookings, &entity.Booking{
		ID: bookingID, CourtID: "64b7f1a2c3d4e5f607182930", UserID: userID,
		CustomerName: "Budi", CustomerPhone: "081234567890",
		Date:      time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
		StartTime: "18:00", EndTime: "20:00", TotalPrice: 200000,
		Status: status, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	_ = courts
	return bookings, bookingID, userID
}

func lifecycleRouter(bookings repository.BookingRepository, userID string, admin bool) *gin.Engine {
	r := gin.New()
	if userID != "" {
		r.Use(func(c *gin.Context) {
			if oid, err := parseHexID(userID); err == nil {
				c.Set("userID", oid)
			}
			if admin {
				c.Set("role", "admin")
			}
			c.Next()
		})
	}
	h := NewBookingHandler(bookings, testutil.NewFakeCourtRepository())
	r.DELETE("/api/bookings/:id", h.CancelBooking)
	r.GET("/api/admin/bookings", h.GetAllBookings)
	r.PATCH("/api/admin/bookings/:id/status", h.UpdateBookingStatus)
	return r
}

func doJSON(t *testing.T, r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	r.ServeHTTP(w, req)
	return w
}

func TestCancelBookingByOwner(t *testing.T) {
	for _, status := range []entity.BookingStatus{entity.BookingStatusPending, entity.BookingStatusConfirmed} {
		bookings, id, userID := seedLifecycleBooking(t, status)
		r := lifecycleRouter(bookings, userID, false)
		w := doJSON(t, r, http.MethodDelete, "/api/bookings/"+id, "")
		if w.Code != http.StatusOK {
			t.Fatalf("status=%s: expected 200, got %d (body: %s)", status, w.Code, w.Body.String())
		}
		if bookings.Bookings[0].Status != entity.BookingStatusCancelled {
			t.Errorf("status=%s: booking not cancelled in store", status)
		}
	}
}

func TestCancelBookingForeignUserGets404(t *testing.T) {
	bookings, id, _ := seedLifecycleBooking(t, entity.BookingStatusPending)
	r := lifecycleRouter(bookings, "64b7f1a2c3d4e5f607182999", false)
	w := doJSON(t, r, http.MethodDelete, "/api/bookings/"+id, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for foreign booking, got %d", w.Code)
	}
	if bookings.Bookings[0].Status != entity.BookingStatusPending {
		t.Error("foreign cancel must not mutate the booking")
	}
}

func TestCancelBookingTerminalStates(t *testing.T) {
	for _, status := range []entity.BookingStatus{entity.BookingStatusCancelled, entity.BookingStatusCompleted} {
		bookings, id, userID := seedLifecycleBooking(t, status)
		r := lifecycleRouter(bookings, userID, false)
		w := doJSON(t, r, http.MethodDelete, "/api/bookings/"+id, "")
		if w.Code != http.StatusConflict {
			t.Fatalf("status=%s: expected 409, got %d", status, w.Code)
		}
	}
}

func TestCancelBookingInvalidID(t *testing.T) {
	bookings, _, userID := seedLifecycleBooking(t, entity.BookingStatusPending)
	r := lifecycleRouter(bookings, userID, false)
	w := doJSON(t, r, http.MethodDelete, "/api/bookings/not-an-id", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAdminStatusTransitions(t *testing.T) {
	cases := []struct {
		from   entity.BookingStatus
		to     string
		expect int
	}{
		{entity.BookingStatusPending, "confirmed", http.StatusOK},
		{entity.BookingStatusPending, "completed", http.StatusOK},
		{entity.BookingStatusConfirmed, "completed", http.StatusOK},
		{entity.BookingStatusConfirmed, "confirmed", http.StatusConflict},
		{entity.BookingStatusCompleted, "confirmed", http.StatusConflict},
		{entity.BookingStatusCancelled, "confirmed", http.StatusConflict},
		{entity.BookingStatusPending, "cancelled", http.StatusBadRequest},
		{entity.BookingStatusPending, "bogus", http.StatusBadRequest},
	}
	for _, tc := range cases {
		bookings, id, userID := seedLifecycleBooking(t, tc.from)
		r := lifecycleRouter(bookings, userID, true)
		w := doJSON(t, r, http.MethodPatch, "/api/admin/bookings/"+id+"/status", `{"status":"`+tc.to+`"}`)
		if w.Code != tc.expect {
			t.Errorf("%s -> %s: expected %d, got %d (body: %s)", tc.from, tc.to, tc.expect, w.Code, w.Body.String())
		}
	}
}

func TestAdminStatusNotFound(t *testing.T) {
	bookings, _, userID := seedLifecycleBooking(t, entity.BookingStatusPending)
	r := lifecycleRouter(bookings, userID, true)
	w := doJSON(t, r, http.MethodPatch, "/api/admin/bookings/64b7f1a2c3d4e5f607182999/status", `{"status":"confirmed"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestAdminGetAllBookingsPagination(t *testing.T) {
	bookings := testutil.NewFakeBookingRepository()
	for i := 0; i < 3; i++ {
		bookings.Bookings = append(bookings.Bookings, &entity.Booking{
			ID: "64b7f1a2c3d4e5f60718294" + string(rune('0'+i)), Status: entity.BookingStatusPending,
		})
	}
	r := lifecycleRouter(bookings, "64b7f1a2c3d4e5f607182931", true)
	w := doJSON(t, r, http.MethodGet, "/api/admin/bookings?limit=2", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body struct {
		Data []*entity.Booking `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Data) != 2 {
		t.Errorf("expected 2 bookings with limit=2, got %d", len(body.Data))
	}
}
