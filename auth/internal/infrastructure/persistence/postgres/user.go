package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth/internal/domain/entity"
	"auth/internal/domain/interfaces"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserPGRepository(db *pgxpool.Pool) interfaces.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Save(ctx context.Context, u *entity.User) error {
	tag, err := r.db.Exec(ctx, `
        INSERT INTO users (id, email, password, role, created_at)
        VALUES ($1,$2,$3,$4,$5)
        ON CONFLICT (email) DO NOTHING
    `, u.ID, u.Email, u.PasswordHash, u.Role, u.CreatedAt)
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return interfaces.ErrUserAlreadyExists
	}

	return nil
}

func (r *UserRepo) GetUser(ctx context.Context, email string) (*entity.User, error) {
	row := r.db.QueryRow(ctx, `
        SELECT id, email, password, role, created_at
        FROM users WHERE email=$1
    `, email)

	var u entity.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, fmt.Errorf("get by email: %w", err)
	}
	return &u, nil
}
