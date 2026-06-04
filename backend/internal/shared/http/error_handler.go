package http

import (
	stdErrors "errors"

	"github.com/gofiber/fiber/v2"

	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
	validatorv10 "github.com/go-playground/validator/v10"
)

func ErrorHandler(
	c *fiber.Ctx,
	err error,
) error {
	var validationErrors validatorv10.ValidationErrors

	if stdErrors.As(err, &validationErrors) {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"error": validationErrors.Error(),
			},
		)
	}

	switch {

	case stdErrors.Is(err, customerrors.ErrUserNotFound):
		return c.Status(
			fiber.StatusNotFound,
		).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)

	case stdErrors.Is(
		err,
		customerrors.ErrUserAlreadyExists,
	):
		return c.Status(
			fiber.StatusConflict,
		).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)

	case stdErrors.Is(
		err,
		customerrors.ErrUnauthorized,
	):
		return c.Status(
			fiber.StatusUnauthorized,
		).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)

	case stdErrors.Is(
		err,
		customerrors.ErrInvalidCredentials,
	):
		return c.Status(
			fiber.StatusUnauthorized,
		).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)

	default:
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}
}
