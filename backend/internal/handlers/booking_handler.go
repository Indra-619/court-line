package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/your-username/book-lapangan/backend/internal/database"
	"github.com/your-username/book-lapangan/backend/internal/models"
	"github.com/your-username/book-lapangan/backend/pkg/pricing"
	"github.com/your-username/book-lapangan/backend/pkg/validate"
)

// nowFn is injectable for tests.
var nowFn = time.Now

// CreateBooking creates a new booking (requires authentication)
func CreateBooking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
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

	// Validate court ID
	courtObjID, err := primitive.ObjectIDFromHex(input.CourtID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	// Check if court exists
	courtCollection := database.Client.Database("booklapangan").Collection("courts")
	var court models.Court
	err = courtCollection.FindOne(ctx, bson.M{"_id": courtObjID}).Decode(&court)
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

	userObjID := userID.(primitive.ObjectID)

	booking := models.Booking{
		ID:            primitive.NewObjectID(),
		CourtID:       courtObjID,
		UserID:        userObjID,
		CustomerName:  input.CustomerName,
		CustomerPhone: input.CustomerPhone,
		Date:          input.Date,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		TotalPrice:    totalPrice,
		Status:        models.BookingStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	collection := database.Client.Database("booklapangan").Collection("bookings")
	_, err = collection.InsertOne(ctx, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": booking})
}

// GetBookings returns bookings for the authenticated user
func GetBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userObjID := userID.(primitive.ObjectID)

	collection := database.Client.Database("booklapangan").Collection("bookings")
	cursor, err := collection.Find(ctx, bson.M{"userId": userObjID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode bookings"})
		return
	}

	if bookings == nil {
		bookings = []models.Booking{}
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
func toBookingsPublicView(bookings []models.Booking) []BookingPublicView {
	views := make([]BookingPublicView, 0, len(bookings))
	for _, b := range bookings {
		views = append(views, BookingPublicView{
			Date:      b.Date,
			StartTime: b.StartTime,
			EndTime:   b.EndTime,
			Status:    string(b.Status),
			ID:        b.ID.Hex(),
		})
	}
	return views
}

// GetBookingsByCourtID returns all bookings for a specific court
func GetBookingsByCourtID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courtIDParam := c.Param("id")
	courtObjID, err := primitive.ObjectIDFromHex(courtIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	collection := database.Client.Database("booklapangan").Collection("bookings")
	cursor, err := collection.Find(ctx, bson.M{"courtId": courtObjID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toBookingsPublicView(bookings)})
}
