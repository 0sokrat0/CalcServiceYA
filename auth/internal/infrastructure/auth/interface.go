package auth

type JWT interface {
	GenerateAccessToken(userID, email, role string) (string, error)
	GenerateRefreshToken(userID, email, role string) (string, error)
	ValidateAccessToken(tokenStr string) (*CustomClaims, error)
	ValidateRefreshToken(tokenStr string) (*CustomClaims, error)
}
