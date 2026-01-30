package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

// Booking represents a court booking
type Booking struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CourtID       primitive.ObjectID `bson:"courtId" json:"courtId"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	CustomerName  string             `bson:"customerName" json:"customerName"`
	CustomerPhone string             `bson:"customerPhone" json:"customerPhone"`
	Date          string             `bson:"date" json:"date"`           // Format: YYYY-MM-DD
	StartTime     string             `bson:"startTime" json:"startTime"` // Format: HH:MM
	EndTime       string             `bson:"endTime" json:"endTime"`     // Format: HH:MM
	TotalPrice    float64            `bson:"totalPrice" json:"totalPrice"`
	Status        BookingStatus      `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// CreateBookingInput represents input for creating a booking
type CreateBookingInput struct {
	CourtID       string `json:"courtId" binding:"required"`
	CustomerName  string `json:"customerName" binding:"required"`
	CustomerPhone string `json:"customerPhone" binding:"required"`
	Date          string `json:"date" binding:"required"`
	StartTime     string `json:"startTime" binding:"required"`
	EndTime       string `json:"endTime" binding:"required"`
}
