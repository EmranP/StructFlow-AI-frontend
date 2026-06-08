package handler

import (
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/dto"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/usecase"

	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUC    usecase.AuthUseCase
	sessionUC usecase.AuthSessionUseCase

	validator *validator.Validator
}

func New(
	authUC usecase.AuthUseCase,
	sessionUC usecase.AuthSessionUseCase,
	validator *validator.Validator,
) *AuthHandler {

	return &AuthHandler{
		authUC:    authUC,
		sessionUC: sessionUC,
		validator: validator,
	}
}

func (h *AuthHandler) Register(
	c *fiber.Ctx,
) error {

	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Validate(&req); err != nil {
		return err
	}

	err := h.authUC.Register(
		c.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		return err
	}

	return c.Status(
		fiber.StatusCreated,
	).JSON(fiber.Map{
		"message": "user created and resend email code",
	})
}

func (h *AuthHandler) Login(
	c *fiber.Ctx,
) error {

	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Validate(&req); err != nil {
		return err
	}

	tokens, err := h.authUC.Login(
		c.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		return err
	}

	c.Cookie(
		&fiber.Cookie{
			Name:     "refreshToken",
			Value:    tokens.Refresh,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/api",
			MaxAge:   60 * 60 * 24 * 30,
		},
	)

	return c.JSON(
		dto.AuthResponse{
			AccessToken: tokens.Access,
		},
	)
}

func (h *AuthHandler) VerifyEmail(
	c *fiber.Ctx,
) error {

	var req dto.VerifyEmailRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Validate(&req); err != nil {
		return err
	}

	tokens, err := h.authUC.VerifyEmail(
		c.Context(),
		req.Email,
		req.Code,
	)

	if err != nil {
		return err
	}

	c.Cookie(
		&fiber.Cookie{
			Name:     "refreshToken",
			Value:    tokens.Refresh,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/api",
			MaxAge:   60 * 60 * 24 * 30,
		},
	)

	return c.JSON(
		fiber.Map{
			"message":     "email verified",
			"accessToken": tokens.Access,
		},
	)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized: user id not found in token",
		})
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id format",
		})
	}

	user, errMe := h.authUC.Me(
		c.Context(),
		userID,
	)
	if errMe != nil {
		return err
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(
		fiber.Map{
			"id": user,
		},
	)
}

// Session
func (h *AuthHandler) Refresh(
	c *fiber.Ctx,
) error {

	refreshToken := c.Cookies(
		"refreshToken",
	)

	if refreshToken == "" {
		return customerrors.ErrUnauthorized
	}

	tokens, err := h.sessionUC.
		Refresh(
			c.Context(),
			refreshToken,
		)

	if err != nil {

		return err
	}

	c.Cookie(
		&fiber.Cookie{
			Name:     "refreshToken",
			Value:    tokens.Refresh,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/api",
			MaxAge:   60 * 60 * 24 * 30,
		},
	)

	return c.JSON(
		fiber.Map{
			"accessToken": tokens.Access,
		},
	)
}

func (h *AuthHandler) Logout(
	c *fiber.Ctx,
) error {

	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized: user id not found in token",
		})
	}

	userID, errParse := uuid.Parse(userIdStr)
	if errParse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id format",
		})
	}

	err := h.sessionUC.Clear(
		c.Context(),
		userID,
	)

	if err != nil {
		return err
	}

	c.Cookie(
		&fiber.Cookie{
			Name:   "refreshToken",
			Value:  "",
			MaxAge: -1,
			Path:   "/api",
		},
	)

	return c.JSON(
		fiber.Map{
			"message": "logout successful",
		},
	)
}
