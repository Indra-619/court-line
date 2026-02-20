package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/your-username/book-lapangan/backend/internal/domain/entity"
)

// Mock Repo
type mockCourtRepo struct {
	findAllFn  func() ([]*entity.Court, error)
	findByIDFn func(id string) (*entity.Court, error)
}

func (m *mockCourtRepo) FindAll(ctx context.Context) ([]*entity.Court, error) { return m.findAllFn() }
func (m *mockCourtRepo) FindByID(ctx context.Context, id string) (*entity.Court, error) {
	return m.findByIDFn(id)
}
func (m *mockCourtRepo) Create(ctx context.Context, court *entity.Court) error { return nil }
func (m *mockCourtRepo) Update(ctx context.Context, court *entity.Court) error { return nil }
func (m *mockCourtRepo) Delete(ctx context.Context, id string) error           { return nil }

func TestCourtUseCase_GetAllCourts(t *testing.T) {
	mockRepo := &mockCourtRepo{
		findAllFn: func() ([]*entity.Court, error) {
			return []*entity.Court{{ID: "1", Name: "Court A"}}, nil
		},
	}
	uc := NewCourtUseCase(mockRepo)

	res := uc.GetAllCourts(context.Background())
	if res.IsFailure() {
		t.Errorf("Expected success, got failure: %v", res.Error)
	}
	if len(res.Value) != 1 {
		t.Errorf("Expected 1 court, got %d", len(res.Value))
	}
}

func TestCourtUseCase_CreateCourt(t *testing.T) {
	mockRepo := &mockCourtRepo{}
	uc := NewCourtUseCase(mockRepo)

	res := uc.CreateCourt(context.Background(), &entity.Court{Name: "New Court"})
	if res.IsFailure() {
		t.Errorf("Expected success, got failure: %v", res.Error)
	}
}

func TestCourtUseCase_UpdateCourt(t *testing.T) {
	mockRepo := &mockCourtRepo{}
	uc := NewCourtUseCase(mockRepo)

	res := uc.UpdateCourt(context.Background(), &entity.Court{ID: "1", Name: "Updated Name"})
	if res.IsFailure() {
		t.Errorf("Expected success, got failure: %v", res.Error)
	}
}

func TestCourtUseCase_DeleteCourt(t *testing.T) {
	mockRepo := &mockCourtRepo{}
	uc := NewCourtUseCase(mockRepo)

	res := uc.DeleteCourt(context.Background(), "1")
	if res.IsFailure() {
		t.Errorf("Expected success, got failure: %v", res.Error)
	}
}

func TestCourtUseCase_Errors(t *testing.T) {
	err := errors.New("db error")
	mockRepo := &mockCourtRepo{
		findAllFn:  func() ([]*entity.Court, error) { return nil, err },
		findByIDFn: func(id string) (*entity.Court, error) { return nil, err },
	}
	uc := NewCourtUseCase(mockRepo)

	if !uc.GetAllCourts(context.Background()).IsFailure() {
		t.Error("Expected error on GetAllCourts")
	}
	if !uc.GetCourtByID(context.Background(), "1").IsFailure() {
		t.Error("Expected error on GetCourtByID")
	}
}
