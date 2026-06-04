package handler

import "github.com/gofiber/fiber/v2"

func Me(
	c *fiber.Ctx,
) error {

	userID := c.Locals("userId")

	return c.JSON(
		fiber.Map{
			"userId": userID,
		},
	)
}
