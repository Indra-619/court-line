package entity

import (
	"time"
)

type Court struct {
	ID           string
	Name         string
	Type         string
	Location     string
	Description  string
	PricePerHour float64
	ImageURL     string
	Facilities   []string
	IsAvailable  bool
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID            string
	CourtID       string
	UserID        string
	CustomerName  string
	CustomerPhone string
	Date          string // YYYY-MM-DD
	StartTime     string // HH:MM
	EndTime       string // HH:MM
	TotalPrice    float64
	Status        BookingStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type User struct {
	ID        string
	GoogleID  string
	Email     string
	Name      string
	Picture   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
