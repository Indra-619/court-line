package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Court represents a sports venue
type Court struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Type         string             `bson:"type" json:"type"` // Futsal, Badminton, Basketball, etc.
	Location     string             `bson:"location" json:"location"`
	Description  string             `bson:"description" json:"description"`
	PricePerHour float64            `bson:"pricePerHour" json:"pricePerHour"`
	ImageURL     string             `bson:"imageUrl" json:"imageUrl"`
	Facilities   []string           `bson:"facilities" json:"facilities"`
	IsAvailable  bool               `bson:"isAvailable" json:"isAvailable"`
}

// CreateCourtInput represents input for creating a court
type CreateCourtInput struct {
	Name         string   `json:"name" binding:"required"`
	Type         string   `json:"type" binding:"required"`
	Location     string   `json:"location" binding:"required"`
	Description  string   `json:"description"`
	PricePerHour float64  `json:"pricePerHour" binding:"required"`
	ImageURL     string   `json:"imageUrl"`
	Facilities   []string `json:"facilities"`
}
