package entity

import (
	"errors"
	"regexp"
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

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func NewUser(email, password, role string) (*User, error) {
	if !emailRegex.MatchString(email) {
		return nil, errors.New("validation error: invalid email format")
	}
	return &User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: password,
		Role:         role,
		CreatedAt:    time.Now(),
	}, nil
}
