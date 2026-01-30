package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/your-username/book-lapangan/backend/internal/database"
	"github.com/your-username/book-lapangan/backend/internal/models"
)

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

	// Calculate total price
	startParts := strings.Split(input.StartTime, ":")
	endParts := strings.Split(input.EndTime, ":")
	startHour, _ := strconv.Atoi(startParts[0])
	endHour, _ := strconv.Atoi(endParts[0])
	hours := float64(endHour - startHour)
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

	if bookings == nil {
		bookings = []models.Booking{}
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}
