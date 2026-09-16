package handlers

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-username/book-lapangan/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToBookingsPublicViewFields(t *testing.T) {
	bookings := []models.Booking{
		{
			ID:            primitive.NewObjectID(),
			CourtID:       primitive.NewObjectID(),
			UserID:        primitive.NewObjectID(),
			CustomerName:  "Budi Santoso",
			CustomerPhone: "081234567890",
			Date:          "2026-09-20",
			StartTime:     "18:00",
			EndTime:       "20:00",
			TotalPrice:    150000,
			Status:        models.BookingStatusPending,
		},
	}

	views := toBookingsPublicView(bookings)
	if len(views) != 1 {
		t.Fatalf("expected 1 view, got %d", len(views))
	}

	view := views[0]
	if view.Date != "2026-09-20" || view.StartTime != "18:00" || view.EndTime != "20:00" {
		t.Errorf("public view lost schedule data: %+v", view)
	}
	if view.Status != string(models.BookingStatusPending) {
		t.Errorf("expected status pending, got %s", view.Status)
	}
	if view.ID != bookings[0].ID.Hex() {
		t.Errorf("expected ID %s, got %s", bookings[0].ID.Hex(), view.ID)
	}
}

func TestBookingPublicViewNoPIIInJSON(t *testing.T) {
	bookings := []models.Booking{
		{
			ID:            primitive.NewObjectID(),
			UserID:        primitive.NewObjectID(),
			CustomerName:  "Budi Santoso",
			CustomerPhone: "081234567890",
			Date:          "2026-09-20",
			StartTime:     "18:00",
			EndTime:       "20:00",
			TotalPrice:    150000,
			Status:        models.BookingStatusConfirmed,
		},
	}

	views := toBookingsPublicView(bookings)
	data, err := json.Marshal(views)
	if err != nil {
		t.Fatalf("failed to marshal public views: %v", err)
	}
	out := string(data)

	for _, forbidden := range []string{
		"customerName", "customerPhone", "userId", "totalPrice",
		"Budi Santoso", "081234567890",
	} {
		if strings.Contains(out, forbidden) {
			t.Errorf("public booking JSON leaks %q: %s", forbidden, out)
		}
	}

	for _, required := range []string{"date", "startTime", "endTime", "status", "id"} {
		if !strings.Contains(out, required) {
			t.Errorf("public booking JSON missing required field %q: %s", required, out)
		}
	}
}

func TestToBookingsPublicViewEmpty(t *testing.T) {
	views := toBookingsPublicView(nil)
	if views == nil {
		t.Error("expected non-nil empty slice for nil input")
	}
	if len(views) != 0 {
		t.Errorf("expected 0 views, got %d", len(views))
	}
}
