package repo

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mrmamongo/go-auth-tg/internal/entity"
	"github.com/mrmamongo/go-auth-tg/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) UserRepo {
	return UserRepo{pg}
}

func (r *UserRepo) GetAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	query, args, err := squirrel.Select("*").From("users").ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.Pool.Query(ctx, query, args...)
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Username, &user.TelegramUsername); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)
	sql, args, err := r.Builder.Select("id", "username", "telegram_username").From("users").Where(squirrel.Eq{"username": username}).Limit(1).ToSql()

	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUserByUsername - r.Builder.Insert: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Username, &user.TelegramUsername)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUserByUsername - r.Pool.QueryRow: %w", err)
	}

	return user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	sql, args, err := r.Builder.Insert("user").Columns("username", "telegram_username").Values(user.Username, user.TelegramUsername).Suffix("RETURNING id").ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Builder.Insert: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Pool.QueryRow: %w", err)
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	sql, args, err := r.Builder.Update("user").Set("username", user.Username).Set("telegram_username", user.TelegramUsername).Where(squirrel.Eq{"id": user.Id}).Suffix("RETURNING id").ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - UpdateTg - r.Builder.Update: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("UserRepo - UpdateTg - r.Pool.QueryRow: %w", err)
	}

	return nil
}
