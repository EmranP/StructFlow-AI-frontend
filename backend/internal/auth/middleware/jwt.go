package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/token"
	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
)

type JWTMiddleware struct {
	tokenService *token.Service
}

func NewJWTMiddleware(
	tokenService *token.Service,
) *JWTMiddleware {
	return &JWTMiddleware{
		tokenService: tokenService,
	}
}

func (m *JWTMiddleware) Protected(
	c *fiber.Ctx,
) error {

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return customerrors.ErrUnauthorized
	}

	parts := strings.Split(
		authHeader,
		" ",
	)

	if len(parts) != 2 {
		return customerrors.ErrUnauthorized
	}

	userID, err := m.tokenService.
		ParseAccessToken(parts[1])

	if err != nil {
		return customerrors.ErrUnauthorized
	}

	c.Locals(
		"userId",
		userID,
	)

	return c.Next()
}
