// Package entity holds the core domain types. JSON tags mirror the
// persistence models so HTTP response bodies stay byte-identical when
// handlers serialize entities directly.
package entity

import (
	"time"
)

type Court struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Location     string   `json:"location"`
	Description  string   `json:"description"`
	PricePerHour float64  `json:"pricePerHour"`
	ImageURL     string   `json:"imageUrl"`
	Facilities   []string `json:"facilities"`
	IsAvailable  bool     `json:"isAvailable"`
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID            string        `json:"id"`
	CourtID       string        `json:"courtId"`
	UserID        string        `json:"userId"`
	CustomerName  string        `json:"customerName"`
	CustomerPhone string        `json:"customerPhone"`
	Date          string        `json:"date"`      // Format: YYYY-MM-DD
	StartTime     string        `json:"startTime"` // Format: HH:MM
	EndTime       string        `json:"endTime"`   // Format: HH:MM
	TotalPrice    float64       `json:"totalPrice"`
	Status        BookingStatus `json:"status"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

type User struct {
	ID        string    `json:"id"`
	GoogleID  string    `json:"googleId"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
