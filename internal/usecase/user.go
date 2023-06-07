package usecase

import (
	"context"
	"github.com/mrmamongo/go-auth-tg/internal/entity"
	"github.com/mrmamongo/go-auth-tg/internal/usecase/repo"
)

type UserUseCase struct {
	repo repo.UserRepo
}

func NewUserUseCase(repo repo.UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (uc *UserUseCase) Create(ctx context.Context, user *entity.User) error {
	return uc.repo.Create(ctx, user)
}

func (uc *UserUseCase) Update(ctx context.Context, user *entity.User) error {
	return uc.repo.Update(ctx, user)
}

func (uc *UserUseCase) GetAll(ctx context.Context) ([]entity.User, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *UserUseCase) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return uc.repo.GetByUsername(ctx, username)
}
