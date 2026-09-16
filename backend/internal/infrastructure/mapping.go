// Package infrastructure contains the MongoDB implementations of the
// domain repository interfaces. Constructors receive an explicit
// *mongo.Client and database name; no package-level state is used.
package infrastructure

import (
	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/models"
)

// courtToEntity maps a bson persistence model to the domain entity.
func courtToEntity(m *models.Court) *entity.Court {
	return &entity.Court{
		ID:           m.ID.Hex(),
		Name:         m.Name,
		Type:         m.Type,
		Location:     m.Location,
		Description:  m.Description,
		PricePerHour: m.PricePerHour,
		ImageURL:     m.ImageURL,
		Facilities:   m.Facilities,
		IsAvailable:  m.IsAvailable,
	}
}

// bookingToEntity maps a bson persistence model to the domain entity.
func bookingToEntity(m *models.Booking) *entity.Booking {
	return &entity.Booking{
		ID:            m.ID.Hex(),
		CourtID:       m.CourtID.Hex(),
		UserID:        m.UserID.Hex(),
		CustomerName:  m.CustomerName,
		CustomerPhone: m.CustomerPhone,
		Date:          m.Date,
		StartTime:     m.StartTime,
		EndTime:       m.EndTime,
		TotalPrice:    m.TotalPrice,
		Status:        entity.BookingStatus(m.Status),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// userToEntity maps a bson persistence model to the domain entity.
func userToEntity(m *models.User) *entity.User {
	return &entity.User{
		ID:        m.ID.Hex(),
		GoogleID:  m.GoogleID,
		Email:     m.Email,
		Name:      m.Name,
		Picture:   m.Picture,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
