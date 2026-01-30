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
)

// GetCourts returns all courts
func GetCourts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Client.Database("booklapangan").Collection("courts")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courts"})
		return
	}
	defer cursor.Close(ctx)

	var courts []models.Court
	if err := cursor.All(ctx, &courts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode courts"})
		return
	}

	if courts == nil {
		courts = []models.Court{}
	}

	c.JSON(http.StatusOK, gin.H{"data": courts})
}

// GetCourtByID returns a single court by ID
func GetCourtByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	collection := database.Client.Database("booklapangan").Collection("courts")

	var court models.Court
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&court)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": court})
}

// CreateCourt creates a new court
func CreateCourt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input models.CreateCourtInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	court := models.Court{
		ID:           primitive.NewObjectID(),
		Name:         input.Name,
		Type:         input.Type,
		Location:     input.Location,
		Description:  input.Description,
		PricePerHour: input.PricePerHour,
		ImageURL:     input.ImageURL,
		Facilities:   input.Facilities,
		IsAvailable:  true,
	}

	collection := database.Client.Database("booklapangan").Collection("courts")
	_, err := collection.InsertOne(ctx, court)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create court"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": court})
}

// UpdateCourt updates an existing court
func UpdateCourt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	var input models.CreateCourtInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.Client.Database("booklapangan").Collection("courts")

	update := bson.M{
		"$set": bson.M{
			"name":         input.Name,
			"type":         input.Type,
			"location":     input.Location,
			"description":  input.Description,
			"pricePerHour": input.PricePerHour,
			"imageUrl":     input.ImageURL,
			"facilities":   input.Facilities,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update court"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Court updated successfully"})
}

// DeleteCourt deletes a court
func DeleteCourt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid court ID"})
		return
	}

	collection := database.Client.Database("booklapangan").Collection("courts")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete court"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Court not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Court deleted successfully"})
}
