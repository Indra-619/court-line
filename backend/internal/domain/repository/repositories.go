// Package repository defines the storage contracts consumed by handlers
// and middleware. Implementations live in internal/infrastructure; the
// rest of the application depends only on these interfaces.
package repository

import (
	"context"
	"errors"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
)

// Sentinel errors returned by repository implementations so callers can
// map storage failures to HTTP responses without knowing the backend.
var (
	// ErrNotFound is returned when a lookup matched no document.
	ErrNotFound = errors.New("not found")
	// ErrInvalidID is returned when an identifier cannot be parsed
	// (e.g. a string that is not a 24-character hex ObjectID).
	ErrInvalidID = errors.New("invalid id")
)

type CourtRepository interface {
	FindAll(ctx context.Context) ([]*entity.Court, error)
	FindByID(ctx context.Context, id string) (*entity.Court, error)
	Create(ctx context.Context, court *entity.Court) error
	Update(ctx context.Context, court *entity.Court) error
	Delete(ctx context.Context, id string) error
}

type BookingRepository interface {
	FindAll(ctx context.Context) ([]*entity.Booking, error)
	FindByID(ctx context.Context, id string) (*entity.Booking, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error)
	// FindByCourtID returns every booking for a court, any status.
	FindByCourtID(ctx context.Context, courtID string) ([]*entity.Booking, error)
	// FindActiveByCourtAndDate returns pending/confirmed bookings for a
	// court on a given date; used for the double-booking overlap check.
	FindActiveByCourtAndDate(ctx context.Context, courtID, date string) ([]*entity.Booking, error)
	Create(ctx context.Context, booking *entity.Booking) error
	Update(ctx context.Context, booking *entity.Booking) error
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
	// UpsertGoogleUser inserts a user for the Google account or refreshes
	// the profile fields of the existing one (role is set only on insert).
	UpsertGoogleUser(ctx context.Context, user *entity.User) error
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
}
