package middleware

import (
	"errors"
	"strings"

	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/http/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("access_token")
		if strings.TrimSpace(tokenStr) == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "missing access_token cookie"})
		}

		claims := new(dto.CustomClaims)
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid or expired JWT"})
		}

		c.Locals("jwt", token)
		return c.Next()
	}
}
