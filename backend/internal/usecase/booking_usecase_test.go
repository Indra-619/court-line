package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/your-username/book-lapangan/backend/internal/domain/entity"
)

type mockBookingRepo struct {
	createFn       func(b *entity.Booking) error
	findAllFn      func() ([]*entity.Booking, error)
	findByUserIDFn func(userID string) ([]*entity.Booking, error)
}

func (m *mockBookingRepo) FindAll(ctx context.Context) ([]*entity.Booking, error) {
	if m.findAllFn != nil {
		return m.findAllFn()
	}
	return nil, nil
}
func (m *mockBookingRepo) FindByID(ctx context.Context, id string) (*entity.Booking, error) {
	return nil, nil
}
func (m *mockBookingRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error) {
	if m.findByUserIDFn != nil {
		return m.findByUserIDFn(userID)
	}
	return nil, nil
}
func (m *mockBookingRepo) Create(ctx context.Context, b *entity.Booking) error { return m.createFn(b) }
func (m *mockBookingRepo) Update(ctx context.Context, b *entity.Booking) error { return nil }

func TestBookingUseCase_CreateBooking_InvalidTime(t *testing.T) {
	mockCourtRepo := &mockCourtRepo{
		findByIDFn: func(id string) (*entity.Court, error) {
			return &entity.Court{ID: id, PricePerHour: 100}, nil
		},
	}
	mockBookingRepo := &mockBookingRepo{}
	uc := NewBookingUseCase(mockBookingRepo, mockCourtRepo)

	// End time before start time
	booking := &entity.Booking{
		CourtID:   "1",
		StartTime: "10:00",
		EndTime:   "09:00",
	}

	res := uc.CreateBooking(context.Background(), booking)
	if res.IsSuccess() {
		t.Error("Expected failure for invalid time range, got success")
	}
	if res.Error.Error() != "end time must be after start time" {
		t.Errorf("Unexpected error message: %v", res.Error)
	}
}

func TestBookingUseCase_GetMyBookings(t *testing.T) {
	mockBookingRepo := &mockBookingRepo{
		createFn: func(b *entity.Booking) error { return nil },
	}
	uc := NewBookingUseCase(mockBookingRepo, nil)

	res := uc.GetMyBookings(context.Background(), "user-1")
	if res.IsFailure() {
		t.Errorf("Expected success, got failure: %v", res.Error)
	}
}

func TestBookingUseCase_GetBookingsByCourt(t *testing.T) {
	mockBookingRepo := &mockBookingRepo{
		createFn: func(b *entity.Booking) error { return nil },
	}
	uc := NewBookingUseCase(mockBookingRepo, nil)

	res := uc.GetBookingsByCourt(context.Background(), "court-1")
	if res.IsFailure() {
		t.Errorf("Expected success, got failure: %v", res.Error)
	}
}

func TestBookingUseCase_CreateBooking_CourtNotFound(t *testing.T) {
	mockCourtRepo := &mockCourtRepo{
		findByIDFn: func(id string) (*entity.Court, error) {
			return nil, nil // Not found
		},
	}
	uc := NewBookingUseCase(nil, mockCourtRepo)

	res := uc.CreateBooking(context.Background(), &entity.Booking{CourtID: "999"})
	if res.IsSuccess() {
		t.Error("Expected failure for non-existent court, got success")
	}
}

func TestBookingUseCase_CreateBooking_InvalidTimeFormat(t *testing.T) {
	mockCourtRepo := &mockCourtRepo{
		findByIDFn: func(id string) (*entity.Court, error) {
			return &entity.Court{ID: id}, nil
		},
	}
	uc := NewBookingUseCase(nil, mockCourtRepo)

	res := uc.CreateBooking(context.Background(), &entity.Booking{
		CourtID:   "1",
		StartTime: "invalid",
		EndTime:   "invalid",
	})
	if res.IsSuccess() {
		t.Error("Expected failure for invalid time format, got success")
	}
}

func TestBookingUseCase_CreateBooking_Success(t *testing.T) {

	mockCourtRepo := &mockCourtRepo{
		findByIDFn: func(id string) (*entity.Court, error) {
			return &entity.Court{ID: id, PricePerHour: 100}, nil
		},
	}
	mockBookingRepo := &mockBookingRepo{
		createFn: func(b *entity.Booking) error { return nil },
	}
	uc := NewBookingUseCase(mockBookingRepo, mockCourtRepo)

	booking := &entity.Booking{
		CourtID:   "1",
		StartTime: "10:00",
		EndTime:   "12:00",
	}

	res := uc.CreateBooking(context.Background(), booking)
	if res.IsFailure() {
		t.Fatalf("Expected success, got failure: %v", res.Error)
	}

	if res.Value.TotalPrice != 200 {
		t.Errorf("Expected TotalPrice 200, got %f", res.Value.TotalPrice)
	}
}

func TestBookingUseCase_GetMyBookings_Error(t *testing.T) {
	mockBookingRepo := &mockBookingRepo{
		findByUserIDFn: func(uid string) ([]*entity.Booking, error) { return nil, errors.New("db error") },
	}
	uc := NewBookingUseCase(mockBookingRepo, nil)

	if !uc.GetMyBookings(context.Background(), "1").IsFailure() {
		t.Error("Expected failure on repo error")
	}
}

func TestBookingUseCase_GetBookingsByCourt_Error(t *testing.T) {
	mockBookingRepo := &mockBookingRepo{
		findAllFn: func() ([]*entity.Booking, error) { return nil, errors.New("db error") },
	}
	uc := NewBookingUseCase(mockBookingRepo, nil)

	if !uc.GetBookingsByCourt(context.Background(), "1").IsFailure() {
		t.Error("Expected failure on repo error")
	}
}
