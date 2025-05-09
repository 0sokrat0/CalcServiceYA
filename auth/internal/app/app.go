package app

import (
	"auth/internal/app/dto"
	"auth/internal/domain/interfaces"
	"auth/internal/infrastructure/auth"
	hashpass "auth/internal/infrastructure/hashPass"

	"context"
	"errors"
)

var (
	ErrEmptyCredentials   = errors.New("email and password must be provided")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Login(ctx context.Context, in dto.LoginRequest) (dto.LoginResponse, error)
	Register(ctx context.Context, in dto.RegisterRequest) (dto.RegisterResponse, error)
}

type UserService struct {
	repo   interfaces.UserRepository
	jwt    auth.JWT
	hasher hashpass.PasswordHasher
}

func NewUserService(
	repo interfaces.UserRepository,
	jwt auth.JWT,
	hasher hashpass.PasswordHasher,
) AuthService {
	return &UserService{repo: repo, jwt: jwt, hasher: hasher}
}
