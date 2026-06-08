package handler

import (
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/dto"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/usecase"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase

	validator *validator.Validator
}

func New(
	authUC usecase.AuthUseCase,
	validator *validator.Validator,
) *AuthHandler {

	return &AuthHandler{
		authUC:    authUC,
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

	token, err := h.authUC.Login(
		c.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		return err
	}

	return c.JSON(
		dto.AuthResponse{
			AccessToken: token,
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

	err := h.authUC.VerifyEmail(
		c.Context(),
		req.Email,
		req.Code,
	)

	if err != nil {
		return err
	}

	return c.JSON(
		fiber.Map{
			"message": "email verified",
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
