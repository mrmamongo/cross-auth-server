// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/mrmamongo/go-auth-tg/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// User -.
	User interface {
		GetAll(context.Context) ([]entity.User, error)
		GetByUsername(context.Context, string) (*entity.User, error)
		GetByTelegram(context.Context, string) (*entity.User, error)
		Create(context.Context, *entity.User) error
		Update(context.Context, *entity.User) error
	}
)
