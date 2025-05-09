package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

func NewUser(email, password, role string) (*User, error) {
	return &User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: password,
		Role:         role,
		CreatedAt:    time.Now(),
	}, nil
}
