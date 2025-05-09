package interfaces

import (
	"auth/internal/domain/entity"
	"context"
	"errors"
)

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetUser(ctx context.Context, email string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}
