package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/book-lapangan/backend/internal/domain/entity"
	"github.com/your-username/book-lapangan/backend/internal/usecase"
)

type BookingHandler struct {
	bookingUseCase usecase.BookingUseCase
}

func NewBookingHandler(u usecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{
		bookingUseCase: u,
	}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	// Re-use logic to get userID from middleware context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		CourtID       string `json:"courtId"`
		CustomerName  string `json:"customerName"`
		CustomerPhone string `json:"customerPhone"`
		Date          string `json:"date"`
		StartTime     string `json:"startTime"`
		EndTime       string `json:"endTime"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking := &entity.Booking{
		CourtID:       input.CourtID,
		UserID:        userID.(string), // Assuming middleware now sets string ID or we convert
		CustomerName:  input.CustomerName,
		CustomerPhone: input.CustomerPhone,
		Date:          input.Date,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
	}

	result := h.bookingUseCase.CreateBooking(c.Request.Context(), booking)
	if result.IsFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": result.Value})
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	result := h.bookingUseCase.GetMyBookings(c.Request.Context(), userID.(string))
	if result.IsFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Value})
}

func (h *BookingHandler) GetBookingsByCourt(c *gin.Context) {
	courtID := c.Param("id")
	result := h.bookingUseCase.GetBookingsByCourt(c.Request.Context(), courtID)
	if result.IsFailure() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result.Value})
}
