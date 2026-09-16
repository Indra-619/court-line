// Package testutil provides in-memory fakes for the domain repository
// interfaces so handler and route tests run without MongoDB.
package testutil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Indra-619/court-line/backend/internal/domain/entity"
	"github.com/Indra-619/court-line/backend/internal/domain/repository"
)

// NewID returns a deterministic 24-hex-character identifier.
func NewID(n int) string {
	return fmt.Sprintf("%024x", n)
}

// FakeCourtRepository is an in-memory CourtRepository.
type FakeCourtRepository struct {
	Courts map[string]*entity.Court
	Err    error // returned by every operation when set
}

func NewFakeCourtRepository() *FakeCourtRepository {
	return &FakeCourtRepository{Courts: make(map[string]*entity.Court)}
}

func (f *FakeCourtRepository) FindAll(ctx context.Context) ([]*entity.Court, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	out := make([]*entity.Court, 0, len(f.Courts))
	for _, c := range f.Courts {
		out = append(out, c)
	}
	return out, nil
}

func (f *FakeCourtRepository) FindByID(ctx context.Context, id string) (*entity.Court, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	court, ok := f.Courts[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return court, nil
}

func (f *FakeCourtRepository) Create(ctx context.Context, court *entity.Court) error {
	if f.Err != nil {
		return f.Err
	}
	if court.ID == "" {
		court.ID = NewID(len(f.Courts) + 1)
	}
	f.Courts[court.ID] = court
	return nil
}

func (f *FakeCourtRepository) Update(ctx context.Context, court *entity.Court) error {
	if f.Err != nil {
		return f.Err
	}
	existing, ok := f.Courts[court.ID]
	if !ok {
		return repository.ErrNotFound
	}
	existing.Name = court.Name
	existing.Type = court.Type
	existing.Location = court.Location
	existing.Description = court.Description
	existing.PricePerHour = court.PricePerHour
	existing.ImageURL = court.ImageURL
	existing.Facilities = court.Facilities
	existing.IsAvailable = court.IsAvailable
	return nil
}

func (f *FakeCourtRepository) Delete(ctx context.Context, id string) error {
	if f.Err != nil {
		return f.Err
	}
	if _, ok := f.Courts[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.Courts, id)
	return nil
}

// FakeBookingRepository is an in-memory BookingRepository.
type FakeBookingRepository struct {
	Bookings []*entity.Booking
	Err      error // returned by every operation when set
}

func NewFakeBookingRepository() *FakeBookingRepository {
	return &FakeBookingRepository{}
}

func (f *FakeBookingRepository) FindAll(ctx context.Context) ([]*entity.Booking, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Bookings, nil
}

func (f *FakeBookingRepository) FindByID(ctx context.Context, id string) (*entity.Booking, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	for _, b := range f.Bookings {
		if b.ID == id {
			return b, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *FakeBookingRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	out := make([]*entity.Booking, 0)
	for _, b := range f.Bookings {
		if b.UserID == userID {
			out = append(out, b)
		}
	}
	return out, nil
}

func (f *FakeBookingRepository) FindByCourtID(ctx context.Context, courtID string) ([]*entity.Booking, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	out := make([]*entity.Booking, 0)
	for _, b := range f.Bookings {
		if b.CourtID == courtID {
			out = append(out, b)
		}
	}
	return out, nil
}

func (f *FakeBookingRepository) FindActiveByCourtAndDate(ctx context.Context, courtID, date string) ([]*entity.Booking, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	out := make([]*entity.Booking, 0)
	for _, b := range f.Bookings {
		if b.CourtID != courtID || b.Date != date {
			continue
		}
		if b.Status == entity.BookingStatusPending || b.Status == entity.BookingStatusConfirmed {
			out = append(out, b)
		}
	}
	return out, nil
}

func (f *FakeBookingRepository) Create(ctx context.Context, booking *entity.Booking) error {
	if f.Err != nil {
		return f.Err
	}
	if booking.ID == "" {
		booking.ID = NewID(len(f.Bookings) + 1)
	}
	f.Bookings = append(f.Bookings, booking)
	return nil
}

func (f *FakeBookingRepository) Update(ctx context.Context, booking *entity.Booking) error {
	if f.Err != nil {
		return f.Err
	}
	for i, b := range f.Bookings {
		if b.ID == booking.ID {
			booking.UpdatedAt = time.Now()
			f.Bookings[i] = booking
			return nil
		}
	}
	return repository.ErrNotFound
}

// FakeUserRepository is an in-memory UserRepository.
type FakeUserRepository struct {
	Users map[string]*entity.User
	Err   error // returned by every operation when set
}

func NewFakeUserRepository() *FakeUserRepository {
	return &FakeUserRepository{Users: make(map[string]*entity.User)}
}

func (f *FakeUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	user, ok := f.Users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (f *FakeUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	for _, u := range f.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *FakeUserRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	for _, u := range f.Users {
		if u.GoogleID == googleID {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *FakeUserRepository) Create(ctx context.Context, user *entity.User) error {
	if f.Err != nil {
		return f.Err
	}
	if user.ID == "" {
		user.ID = NewID(len(f.Users) + 1)
	}
	f.Users[user.ID] = user
	return nil
}

func (f *FakeUserRepository) Update(ctx context.Context, user *entity.User) error {
	if f.Err != nil {
		return f.Err
	}
	if _, ok := f.Users[user.ID]; !ok {
		return repository.ErrNotFound
	}
	user.UpdatedAt = time.Now()
	f.Users[user.ID] = user
	return nil
}

// FakeRefreshTokenRepository is an in-memory RefreshTokenRepository.
// FindByHash applies the same validity filter as the Mongo
// implementation: expired or revoked records are treated as unknown.
// FindAndDeleteByHash deletes under a lock so concurrent refreshes
// observe the same single-winner rotation semantics as Mongo's
// FindOneAndDelete.
type FakeRefreshTokenRepository struct {
	mu     sync.Mutex
	Tokens map[string]*entity.RefreshToken // keyed by token hash
	Err    error                           // returned by every operation when set
}

func NewFakeRefreshTokenRepository() *FakeRefreshTokenRepository {
	return &FakeRefreshTokenRepository{Tokens: make(map[string]*entity.RefreshToken)}
}

func (f *FakeRefreshTokenRepository) Create(ctx context.Context, record *entity.RefreshToken) error {
	if f.Err != nil {
		return f.Err
	}
	if record.ID == "" {
		record.ID = NewID(len(f.Tokens) + 1)
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}
	f.Tokens[record.TokenHash] = record
	return nil
}

func (f *FakeRefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	record, ok := f.Tokens[hash]
	if !ok || record.Revoked || time.Now().After(record.ExpiresAt) {
		return nil, repository.ErrNotFound
	}
	return record, nil
}

func (f *FakeRefreshTokenRepository) FindAndDeleteByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Err != nil {
		return nil, f.Err
	}
	record, ok := f.Tokens[hash]
	if !ok || record.Revoked || time.Now().After(record.ExpiresAt) {
		return nil, repository.ErrNotFound
	}
	delete(f.Tokens, hash)
	return record, nil
}

func (f *FakeRefreshTokenRepository) DeleteByHash(ctx context.Context, hash string) error {
	if f.Err != nil {
		return f.Err
	}
	if _, ok := f.Tokens[hash]; !ok {
		return repository.ErrNotFound
	}
	delete(f.Tokens, hash)
	return nil
}

func (f *FakeRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	if f.Err != nil {
		return f.Err
	}
	for hash, record := range f.Tokens {
		if record.UserID == userID {
			delete(f.Tokens, hash)
		}
	}
	return nil
}

// FakeRevokedTokenRepository is an in-memory RevokedTokenRepository.
// Entries past their expiry are ignored, mirroring the Mongo TTL index.
type FakeRevokedTokenRepository struct {
	JTIs map[string]time.Time // jti -> expiry of the blacklisted JWT
	Err  error                // returned by every operation when set
}

func NewFakeRevokedTokenRepository() *FakeRevokedTokenRepository {
	return &FakeRevokedTokenRepository{JTIs: make(map[string]time.Time)}
}

func (f *FakeRevokedTokenRepository) Create(ctx context.Context, record *entity.RevokedToken) error {
	if f.Err != nil {
		return f.Err
	}
	f.JTIs[record.JTI] = record.ExpiresAt
	return nil
}

func (f *FakeRevokedTokenRepository) Exists(ctx context.Context, jti string) (bool, error) {
	if f.Err != nil {
		return false, f.Err
	}
	expiresAt, ok := f.JTIs[jti]
	if !ok || time.Now().After(expiresAt) {
		return false, nil
	}
	return true, nil
}

// UpsertGoogleUser mirrors the OAuth callback semantics: refresh the
// profile of an existing user or create one with the default role.
func (f *FakeUserRepository) UpsertGoogleUser(ctx context.Context, user *entity.User) error {
	if f.Err != nil {
		return f.Err
	}
	if existing, err := f.FindByGoogleID(ctx, user.GoogleID); err == nil {
		existing.Email = user.Email
		existing.Name = user.Name
		existing.Picture = user.Picture
		existing.UpdatedAt = time.Now()
		f.Users[existing.ID] = existing
		return nil
	}
	if user.Role == "" {
		user.Role = "user"
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	return f.Create(ctx, user)
}
