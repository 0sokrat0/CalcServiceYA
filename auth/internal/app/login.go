package app

import (
	"auth/internal/app/dto"
	"context"
)

func (s *UserService) Login(ctx context.Context, in dto.LoginRequest) (dto.LoginResponse, error) {
	if in.Email == "" || in.Password == "" {
		return dto.LoginResponse{}, ErrEmptyCredentials
	}

	user, err := s.repo.GetUser(ctx, in.Email)
	if err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}
	if !s.hasher.Compare(user.PasswordHash, in.Password) {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	access, err := s.jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	refresh, err := s.jwt.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Access:  access,
		Refresh: refresh,
	}, nil
}
