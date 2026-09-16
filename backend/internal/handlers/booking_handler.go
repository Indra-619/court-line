package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/domain/repository"
	"github.com/Indra-619/court-line/backend/internal/models"
	"github.com/Indra-619/court-line/backend/pkg/pricing"
	"github.com/Indra-619/court-line/backend/pkg/validate"
)

// nowFn is injectable for tests.
var nowFn = time.Now

// BookingHandler serves the booking endpoints on top of injected
// repositories; it never touches the database driver directly.
type BookingHandler struct {
	bookings repository.BookingRepository
	courts   repository.CourtRepository
}

// NewBookingHandler builds a BookingHandler. The court repository is
// needed to resolve the price and existence of the booked court.
func NewBookingHandler(bookings repository.BookingRepository, courts repository.CourtRepository) *BookingHandler {
	return &BookingHandler{bookings: bookings, courts: courts}
}

// userIDHex extracts the authenticated user's ID (set by the auth
// middleware as a primitive.ObjectID) in hex form.
func userIDHex(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	userObjID, ok := userID.(primitive.ObjectID)
	if !ok {
		return "", false
	}
	return userObjID.Hex(), true
}

// CreateBooking creates a new booking (requires authentication)
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user ID from context (set by auth middleware)
	userID, ok := userIDHex(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input models.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate date/time input format and business rules
	if !validate.ValidateDate(input.Date) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected YYYY-MM-DD"})
		return
	}
	if validate.IsPastDate(input.Date, nowFn()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking date cannot be in the past"})
		return
	}
	if !validate.ValidateClock(input.StartTime) || !validate.ValidateClock(input.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format, expected HH:MM"})
		return
	}
	startMinutes, err := validate.ToMinutes(input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endMinutes, err := validate.ToMinutes(input.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if endMinutes <= startMinutes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}

	// Check if court exists (invalid IDs are rejected up front)
	if _, err := parseHexID(input.CourtID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}
	court, err := h.courts.FindByID(ctx, input.CourtID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
		return
	}

	// Calculate total price from exact minute-based duration
	hours, err := pricing.CalculateHours(input.StartTime, input.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	totalPrice := hours * court.PricePerHour

	// Prevent double-booking: check existing active bookings for the same
	// court and date. NOTE: this is a check-then-insert with a race window
	// on a standalone Mongo (no transactions available); accepted for this
	// project's scale.
	existing, err := h.bookings.FindActiveByCourtAndDate(ctx, input.CourtID, input.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing bookings"})
		return
	}

	for _, b := range existing {
		conflicts, err := pricing.Overlaps(input.StartTime, input.EndTime, b.StartTime, b.EndTime)
		if err != nil {
			// Skip bookings with malformed stored times rather than failing.
			continue
		}
		if conflicts {
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Time slot conflicts with an existing booking (%s-%s)", b.StartTime, b.EndTime),
				"conflict": gin.H{
					"startTime": b.StartTime,
					"endTime":   b.EndTime,
				},
			})
			return
		}
	}

	now := time.Now()
	booking := &entity.Booking{
		CourtID:       input.CourtID,
		UserID:        userID,
		CustomerName:  input.CustomerName,
		CustomerPhone: input.CustomerPhone,
		Date:          input.Date,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		TotalPrice:    totalPrice,
		Status:        entity.BookingStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.bookings.Create(ctx, booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": booking})
}

// GetBookings returns bookings for the authenticated user
func (h *BookingHandler) GetBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, ok := userIDHex(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bookings, err := h.bookings.FindByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	if bookings == nil {
		bookings = []*entity.Booking{}
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

// BookingPublicView exposes only non-sensitive booking fields to
// unauthenticated users browsing a court's schedule.
type BookingPublicView struct {
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Status    string `json:"status"`
	ID        string `json:"id"`
}

// toBookingsPublicView strips customer PII from bookings.
func toBookingsPublicView(bookings []*entity.Booking) []BookingPublicView {
	views := make([]BookingPublicView, 0, len(bookings))
	for _, b := range bookings {
		views = append(views, BookingPublicView{
			Date:      b.Date,
			StartTime: b.StartTime,
			EndTime:   b.EndTime,
			Status:    string(b.Status),
			ID:        b.ID,
		})
	}
	return views
}

// GetBookingsByCourtID returns all bookings for a specific court
func (h *BookingHandler) GetBookingsByCourtID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courtIDParam := c.Param("id")
	if _, err := parseHexID(courtIDParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	bookings, err := h.bookings.FindByCourtID(ctx, courtIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toBookingsPublicView(bookings)})
}
