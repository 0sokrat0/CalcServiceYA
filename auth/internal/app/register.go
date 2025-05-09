package app

import (
	"auth/internal/app/dto"
	"auth/internal/domain/entity"
	"context"
)

func (s *UserService) Register(ctx context.Context, in dto.RegisterRequest) (dto.RegisterResponse, error) {
	if in.Email == "" || in.Password == "" {
		return dto.RegisterResponse{}, ErrEmptyCredentials
	}

	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return dto.RegisterResponse{}, err
	}
	u, err := entity.NewUser(in.Email, hash, "user")
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	if err := s.repo.Save(ctx, u); err != nil {
		return dto.RegisterResponse{}, err
	}

	// TODO: отправить подтверждение на e‑mail

	return dto.RegisterResponse{
		UserID: u.ID,
	}, nil
}
