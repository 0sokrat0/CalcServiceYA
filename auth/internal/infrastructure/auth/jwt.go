package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ClaimTokenType = "token_type"
	RefreshType    = "refresh"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrUnexpectedMethod = errors.New("unexpected signing method")
	ErrInvalidTokenType = errors.New("invalid token type")
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type,omitempty"`
}

type JWTService struct {
	appName         string
	secret          []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewJWTService(secret string, accessDur, refreshDur time.Duration) JWT {
	return &JWTService{
		secret:          []byte(secret),
		accessDuration:  accessDur,
		refreshDuration: refreshDur,
	}
}

func (j *JWTService) GenerateAccessToken(userID, email, role string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    j.appName,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessDuration)),
		},
		Email: email,
		Role:  role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTService) GenerateRefreshToken(userID, email, role string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    j.appName,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshDuration)),
		},
		Email:     email,
		Role:      role,
		TokenType: RefreshType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTService) ValidateAccessToken(tok string) (*CustomClaims, error) {
	return j.validate(tok, "")
}

func (j *JWTService) ValidateRefreshToken(tok string) (*CustomClaims, error) {
	return j.validate(tok, RefreshType)
}

func (j *JWTService) validate(tokenStr, wantType string) (*CustomClaims, error) {
	keyFn := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedMethod
		}
		return j.secret, nil
	}
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, keyFn)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if wantType != "" && claims.TokenType != wantType {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}
