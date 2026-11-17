package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/domain/repository"
)

// transitionInput is the admin status-change payload. Only confirmed
// and completed are reachable through the admin endpoint; cancellation
// is the user-facing DELETE flow.
type transitionInput struct {
	Status string `json:"status" binding:"required"`
}

// findBookingByID parses the :id param and loads the booking, writing
// the appropriate error response itself. Returns ok=false when the
// caller should stop.
func (h *BookingHandler) findBookingByID(c *gin.Context) (*entity.Booking, string, bool) {
	idParam := c.Param("id")
	if _, err := parseHexID(idParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return nil, "", false
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	booking, err := h.bookings.FindByID(ctx, idParam)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
			return nil, "", false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booking"})
		return nil, "", false
	}
	return booking, idParam, true
}

// CancelBooking lets a booking's owner cancel it while it is still
// pending or confirmed. Completed bookings are immutable history.
// DELETE /api/bookings/:id
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID, ok := userIDHex(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	booking, idParam, ok := h.findBookingByID(c)
	if !ok {
		return
	}
	if booking.UserID != userID {
		// Do not reveal that the booking exists for another user.
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	switch booking.Status {
	case entity.BookingStatusPending, entity.BookingStatusConfirmed:
		// cancellable
	case entity.BookingStatusCancelled:
		c.JSON(http.StatusConflict, gin.H{"error": "Booking is already cancelled"})
		return
	default:
		c.JSON(http.StatusConflict, gin.H{"error": "Completed bookings cannot be cancelled"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	if err := h.bookings.UpdateStatus(ctx, idParam, entity.BookingStatusCancelled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}
	booking.Status = entity.BookingStatusCancelled
	c.JSON(http.StatusOK, gin.H{"data": booking, "message": "Booking cancelled"})
}

// UpdateBookingStatus is the admin lifecycle transition:
// pending -> confirmed -> completed. Cancelled bookings are terminal.
// PATCH /api/admin/bookings/:id/status
func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	booking, idParam, ok := h.findBookingByID(c)
	if !ok {
		return
	}

	var input transitionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status payload"})
		return
	}

	var target entity.BookingStatus
	switch input.Status {
	case string(entity.BookingStatusConfirmed):
		target = entity.BookingStatusConfirmed
	case string(entity.BookingStatusCompleted):
		target = entity.BookingStatusCompleted
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be confirmed or completed"})
		return
	}

	allowed := false
	switch booking.Status {
	case entity.BookingStatusPending:
		allowed = target == entity.BookingStatusConfirmed || target == entity.BookingStatusCompleted
	case entity.BookingStatusConfirmed:
		allowed = target == entity.BookingStatusCompleted
	}
	if !allowed {
		c.JSON(http.StatusConflict, gin.H{"error": "Invalid status transition from " + string(booking.Status)})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	if err := h.bookings.UpdateStatus(ctx, idParam, target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}
	booking.Status = target
	c.JSON(http.StatusOK, gin.H{"data": booking})
}

// GetAllBookings lists every booking for the admin dashboard,
// paginated like the other list endpoints.
// GET /api/admin/bookings
func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	bookings, err := h.bookings.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	if bookings == nil {
		bookings = []*entity.Booking{}
	}

	limit, offset := listPagination(c)
	c.JSON(http.StatusOK, gin.H{"data": pageSlice(bookings, limit, offset)})
}
