package usecase

import (
	"context"

	"boilerplate/services/test/domain"
)

type useCase struct {
	repo domain.TestRepository
}

func NewTestUseCase(repo domain.TestRepository) domain.TestUseCase {
	return &useCase{repo: repo}
}

func (u *useCase) Create(ctx context.Context, test *domain.Test) (*domain.Test, error) {
	return u.repo.Create(ctx, test)
}

func (u *useCase) GetByID(ctx context.Context, id string) (*domain.Test, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) List(ctx context.Context) ([]domain.Test, error) {
	return u.repo.List(ctx)
}

func (u *useCase) Update(ctx context.Context, id string, test *domain.Test) (*domain.Test, error) {
	return u.repo.Update(ctx, id, test)
}

func (u *useCase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
