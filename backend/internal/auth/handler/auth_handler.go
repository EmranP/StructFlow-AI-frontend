package handler

import (
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/dto"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/usecase"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v2"
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
		"message": "user created",
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
