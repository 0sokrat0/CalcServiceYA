package app

import (
	"auth/internal/app/dto"
	"auth/internal/domain/entity"
	"auth/internal/infrastructure/auth"
	"context"
	"errors"
	"testing"
)

type mockUserRepo struct {
	getUserFunc func(ctx context.Context, email string) (*entity.User, error)
	saveFunc    func(ctx context.Context, user *entity.User) error
}

func (m *mockUserRepo) GetUser(ctx context.Context, email string) (*entity.User, error) {
	return m.getUserFunc(ctx, email)
}

func (m *mockUserRepo) Save(ctx context.Context, user *entity.User) error {
	return m.saveFunc(ctx, user)
}

type mockHasher struct {
	hashFunc    func(password string) (string, error)
	compareFunc func(hash, password string) bool
}

func (m *mockHasher) Hash(password string) (string, error) {
	return m.hashFunc(password)
}

func (m *mockHasher) Compare(hash, password string) bool {
	return m.compareFunc(hash, password)
}

type CustomClaims struct {
	UserID string
	Email  string
	Role   string
}

// Исправляем mockJWT согласно интерфейсу
type mockJWT struct {
	generateAccessFunc  func(userID, email, role string) (string, error)
	generateRefreshFunc func(userID, email, role string) (string, error)
	validateAccessFunc  func(token string) (*auth.CustomClaims, error)
	validateRefreshFunc func(token string) (*auth.CustomClaims, error)
}

func (m *mockJWT) GenerateAccessToken(userID, email, role string) (string, error) {
	return m.generateAccessFunc(userID, email, role)
}

func (m *mockJWT) GenerateRefreshToken(userID, email, role string) (string, error) {
	return m.generateRefreshFunc(userID, email, role)
}

func (m *mockJWT) ValidateAccessToken(token string) (*auth.CustomClaims, error) {
	return m.validateAccessFunc(token)
}

func (m *mockJWT) ValidateRefreshToken(token string) (*auth.CustomClaims, error) {
	return m.validateRefreshFunc(token)
}

func TestUserService_Login(t *testing.T) {
	emptyRequest := dto.LoginRequest{Email: "", Password: ""}
	validRequest := dto.LoginRequest{Email: "test@test.com", Password: "password"}

	user := entity.User{
		ID:           "123",
		Email:        "test@test.com",
		PasswordHash: "hash",
		Role:         "user",
	}

	tests := []struct {
		name        string
		repo        *mockUserRepo
		hasher      *mockHasher
		jwt         *mockJWT
		request     dto.LoginRequest
		wantError   error
		wantTokens  bool
		description string
	}{
		{
			name:        "empty credentials",
			repo:        &mockUserRepo{},
			hasher:      &mockHasher{},
			jwt:         &mockJWT{},
			request:     emptyRequest,
			wantError:   ErrEmptyCredentials,
			description: "should return error when credentials are empty",
		},
		{
			name: "user not found",
			repo: &mockUserRepo{
				getUserFunc: func(ctx context.Context, email string) (*entity.User, error) {
					return nil, errors.New("not found")
				},
			},
			hasher:      &mockHasher{},
			jwt:         &mockJWT{},
			request:     validRequest,
			wantError:   ErrInvalidCredentials,
			description: "should return error when user not found",
		},
		{
			name: "invalid password",
			repo: &mockUserRepo{
				getUserFunc: func(ctx context.Context, email string) (*entity.User, error) {
					return &user, errors.New("not found")
				},
			},
			hasher: &mockHasher{
				compareFunc: func(hash, password string) bool {
					return false
				},
			},
			jwt:         &mockJWT{},
			request:     validRequest,
			wantError:   ErrInvalidCredentials,
			description: "should return error when password is invalid",
		},
		{
			name: "jwt generation error",
			repo: &mockUserRepo{
				getUserFunc: func(ctx context.Context, email string) (*entity.User, error) {
					// Возвращаем пользователя без ошибки
					return &user, nil
				},
			},
			hasher: &mockHasher{
				compareFunc: func(hash, password string) bool {
					return true // Пароль верный
				},
			},
			jwt: &mockJWT{
				generateAccessFunc: func(userID, email, role string) (string, error) {
					return "", errors.New("jwt error")
				},
			},
			request:     validRequest,
			wantError:   errors.New("jwt error"),
			description: "should return error when token generation fails",
		},
		{
			name: "successful login",
			repo: &mockUserRepo{
				getUserFunc: func(ctx context.Context, email string) (*entity.User, error) {
					return &user, nil
				},
			},
			hasher: &mockHasher{
				compareFunc: func(hash, password string) bool {
					return true // Пароль верный
				},
			},
			jwt: &mockJWT{
				generateAccessFunc: func(userID, email, role string) (string, error) {
					return "access_token", nil
				},
				generateRefreshFunc: func(userID, email, role string) (string, error) {
					return "refresh_token", nil
				},
			},
			request:     validRequest,
			wantTokens:  true,
			description: "should return valid tokens on success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:   tt.repo,
				hasher: tt.hasher,
				jwt:    tt.jwt,
			}

			resp, err := service.Login(context.Background(), tt.request)

			if tt.wantError != nil {
				if err == nil || err.Error() != tt.wantError.Error() {
					t.Errorf("%s\nwant error: %v\ngot: %v", tt.description, tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s\nunexpected error: %v", tt.description, err)
				return
			}

			if tt.wantTokens && (resp.Access == "" || resp.Refresh == "") {
				t.Errorf("%s\nexpected tokens, got empty", tt.description)
			}
		})
	}
}

func TestUserService_Register(t *testing.T) {
	emptyRequest := dto.RegisterRequest{Email: "", Password: ""}
	validRequest := dto.RegisterRequest{Email: "test@test.com", Password: "password"}

	tests := []struct {
		name        string
		repo        *mockUserRepo
		hasher      *mockHasher
		request     dto.RegisterRequest
		wantError   error
		wantUserID  bool
		description string
	}{
		{
			name:        "empty credentials",
			repo:        &mockUserRepo{},
			hasher:      &mockHasher{},
			request:     emptyRequest,
			wantError:   ErrEmptyCredentials,
			description: "should return error when credentials are empty",
		},
		{
			name: "hashing error",
			hasher: &mockHasher{
				hashFunc: func(password string) (string, error) {
					return "", errors.New("hashing error")
				},
			},
			repo:        &mockUserRepo{},
			request:     validRequest,
			wantError:   errors.New("hashing error"),
			description: "should return error when hashing fails",
		},
		{
			name: "invalid user data",
			hasher: &mockHasher{
				hashFunc: func(password string) (string, error) {
					return "hash", nil
				},
			},
			repo: &mockUserRepo{
				saveFunc: func(ctx context.Context, user *entity.User) error {
					t.Error("Save should not be called when user creation fails")
					return nil
				},
			},
			request:     dto.RegisterRequest{Email: "invalid-email", Password: "pass"},
			wantError:   errors.New("validation error: invalid email format"),
			description: "should return error when user data is invalid",
		},
		{
			name: "save error",
			hasher: &mockHasher{
				hashFunc: func(password string) (string, error) {
					return "hash", nil
				},
			},
			repo: &mockUserRepo{
				saveFunc: func(ctx context.Context, user *entity.User) error { // Исправлен тип параметра
					return errors.New("save error")
				},
			},
			request:     validRequest,
			wantError:   errors.New("save error"),
			description: "should return error when save fails",
		},
		{
			name: "successful registration",
			hasher: &mockHasher{
				hashFunc: func(password string) (string, error) {
					return "hash", nil
				},
			},
			repo: &mockUserRepo{
				saveFunc: func(ctx context.Context, user *entity.User) error {
					if user.Email != validRequest.Email || user.PasswordHash != "hash" {
						t.Error("invalid user data")
					}
					return nil
				},
			},
			request:     validRequest,
			wantUserID:  true,
			description: "should return user ID on success",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:   tt.repo,
				hasher: tt.hasher,
				jwt:    &mockJWT{},
			}

			resp, err := service.Register(context.Background(), tt.request)

			if tt.wantError != nil {
				if err == nil {
					t.Errorf("%s\nexpected error: %v\ngot nil", tt.description, tt.wantError)
					return
				}
				if err.Error() != tt.wantError.Error() {
					t.Errorf("%s\nwant error: %v\ngot: %v", tt.description, tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s\nunexpected error: %v", tt.description, err)
				return
			}

			if tt.wantUserID && resp.UserID == "" {
				t.Errorf("%s\nexpected user ID, got empty", tt.description)
			}
		})
	}
}
