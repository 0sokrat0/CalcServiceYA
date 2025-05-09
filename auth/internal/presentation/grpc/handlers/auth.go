package handlers

import (
	"context"

	"auth/internal/app"
	"auth/internal/app/dto"
	pb "auth/pkg/gen/api"
)

type AuthHandler struct {
	uc app.AuthService
	pb.UnimplementedAuthServer
}

func NewAuthHandler(uc app.AuthService) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	in := dto.LoginRequest{Email: req.Email, Password: req.Password}
	out, err := h.uc.Login(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Access:  out.Access,
		Refresh: out.Refresh,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	in := dto.RegisterRequest{Email: req.Email, Password: req.Password}
	out, err := h.uc.Register(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		UserId: out.UserID,
	}, nil
}
