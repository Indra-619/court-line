package repository

import (
	"context"

	"github.com/your-username/book-lapangan/backend/internal/domain/entity"
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
	Create(ctx context.Context, booking *entity.Booking) error
	Update(ctx context.Context, booking *entity.Booking) error
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
}
