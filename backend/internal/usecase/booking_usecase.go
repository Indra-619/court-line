package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/your-username/book-lapangan/backend/internal/domain/entity"
	"github.com/your-username/book-lapangan/backend/internal/domain/repository"
	"github.com/your-username/book-lapangan/backend/pkg/result"
)

type BookingUseCase interface {
	CreateBooking(ctx context.Context, booking *entity.Booking) result.Result[*entity.Booking]
	GetMyBookings(ctx context.Context, userID string) result.Result[[]*entity.Booking]
	GetBookingsByCourt(ctx context.Context, courtID string) result.Result[[]*entity.Booking]
}

type bookingUseCase struct {
	bookingRepo repository.BookingRepository
	courtRepo   repository.CourtRepository
}

func NewBookingUseCase(bookingRepo repository.BookingRepository, courtRepo repository.CourtRepository) BookingUseCase {
	return &bookingUseCase{
		bookingRepo: bookingRepo,
		courtRepo:   courtRepo,
	}
}

func (u *bookingUseCase) CreateBooking(ctx context.Context, booking *entity.Booking) result.Result[*entity.Booking] {
	// 1. Verify court exists
	court, err := u.courtRepo.FindByID(ctx, booking.CourtID)
	if err != nil {
		return result.Failure[*entity.Booking](err)
	}
	if court == nil {
		return result.Failure[*entity.Booking](errors.New("court not found"))
	}

	// 2. Calculate price (logic moved from handler)
	startParts := strings.Split(booking.StartTime, ":")
	endParts := strings.Split(booking.EndTime, ":")
	if len(startParts) != 2 || len(endParts) != 2 {
		return result.Failure[*entity.Booking](errors.New("invalid time format"))
	}

	startHour, _ := strconv.Atoi(startParts[0])
	endHour, _ := strconv.Atoi(endParts[0])
	hours := float64(endHour - startHour)
	if hours <= 0 {
		return result.Failure[*entity.Booking](errors.New("end time must be after start time"))
	}

	booking.TotalPrice = hours * court.PricePerHour
	booking.Status = entity.BookingStatusPending
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	// 3. Save booking
	err = u.bookingRepo.Create(ctx, booking)
	if err != nil {
		return result.Failure[*entity.Booking](err)
	}

	return result.Success(booking)
}

func (u *bookingUseCase) GetMyBookings(ctx context.Context, userID string) result.Result[[]*entity.Booking] {
	bookings, err := u.bookingRepo.FindByUserID(ctx, userID)
	if err != nil {
		return result.Failure[[]*entity.Booking](err)
	}
	return result.Success(bookings)
}

func (u *bookingUseCase) GetBookingsByCourt(ctx context.Context, courtID string) result.Result[[]*entity.Booking] {
	bookings, err := u.bookingRepo.FindAll(ctx) // This should probably be specialized in repo
	if err != nil {
		return result.Failure[[]*entity.Booking](err)
	}

	// Filtering logic if repo doesn't support specific query yet
	var filtered []*entity.Booking
	for _, b := range bookings {
		if b.CourtID == courtID {
			filtered = append(filtered, b)
		}
	}
	return result.Success(filtered)
}
