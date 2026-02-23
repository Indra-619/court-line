package usecase

import (
	"context"

	"github.com/your-username/book-lapangan/backend/internal/domain/entity"
	"github.com/your-username/book-lapangan/backend/internal/domain/repository"
	"github.com/your-username/book-lapangan/backend/pkg/result"
)

type CourtUseCase interface {
	GetAllCourts(ctx context.Context) result.Result[[]*entity.Court]
	GetCourtByID(ctx context.Context, id string) result.Result[*entity.Court]
	CreateCourt(ctx context.Context, court *entity.Court) result.Result[*entity.Court]
	UpdateCourt(ctx context.Context, court *entity.Court) result.Result[string]
	DeleteCourt(ctx context.Context, id string) result.Result[string]
}

type courtUseCase struct {
	courtRepo repository.CourtRepository
}

func NewCourtUseCase(repo repository.CourtRepository) CourtUseCase {
	return &courtUseCase{
		courtRepo: repo,
	}
}

func (u *courtUseCase) GetAllCourts(ctx context.Context) result.Result[[]*entity.Court] {
	courts, err := u.courtRepo.FindAll(ctx)
	if err != nil {
		return result.Failure[[]*entity.Court](err)
	}
	return result.Success(courts)
}

func (u *courtUseCase) GetCourtByID(ctx context.Context, id string) result.Result[*entity.Court] {
	court, err := u.courtRepo.FindByID(ctx, id)
	if err != nil {
		return result.Failure[*entity.Court](err)
	}
	return result.Success(court)
}

func (u *courtUseCase) CreateCourt(ctx context.Context, court *entity.Court) result.Result[*entity.Court] {
	err := u.courtRepo.Create(ctx, court)
	if err != nil {
		return result.Failure[*entity.Court](err)
	}
	return result.Success(court)
}

func (u *courtUseCase) UpdateCourt(ctx context.Context, court *entity.Court) result.Result[string] {
	err := u.courtRepo.Update(ctx, court)
	if err != nil {
		return result.Failure[string](err)
	}
	return result.Success("Court updated successfully")
}

func (u *courtUseCase) DeleteCourt(ctx context.Context, id string) result.Result[string] {
	err := u.courtRepo.Delete(ctx, id)
	if err != nil {
		return result.Failure[string](err)
	}
	return result.Success("Court deleted successfully")
}
