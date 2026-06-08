package handler

import (
	"strconv"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/project/dto"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/project/usecase"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectUC usecase.ProjectUseCase

	validator *validator.Validator
}

func New(
	projectUC usecase.ProjectUseCase,
	validator *validator.Validator,
) *ProjectHandler {
	return &ProjectHandler{
		projectUC: projectUC,
		validator: validator,
	}
}

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	var req dto.ProjectRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Validate(&req); err != nil {
		return err
	}

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

	projectId, err := h.projectUC.Add(
		c.Context(),
		userID,
		&req,
	)
	if err != nil {

		return err
	}

	return c.Status(
		fiber.StatusCreated,
	).JSON(fiber.Map{
		"message": "project created",
		"id":      projectId,
	})
}

func (h *ProjectHandler) GetById(c *fiber.Ctx) error {
	projectId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}

	project, err := h.projectUC.FindByID(c.Context(), projectId)
	if err != nil {
		return err
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(project)
}

func (h *ProjectHandler) GetAll(c *fiber.Ctx) error {

	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	limitParam := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized: user id not found in token",
		})
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	projects, totalCount, err := h.projectUC.FindByUserID(
		c.Context(),
		userID,
		page,
		limit,
	)
	if err != nil {
		return err
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(
		fiber.Map{
			"projects":   projects,
			"page":       page,
			"limit":      limit,
			"totalCount": totalCount,
		},
	)
}

func (h *ProjectHandler) Edit(
	c *fiber.Ctx,
) error {
	projectIdStr := c.Params("id")

	if projectIdStr == "" {
		return c.Status(
			fiber.StatusNotFound,
		).JSON(fiber.Map{
			"message": "Project ID is empty",
		})
	}

	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		return err
	}

	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized: user id not found in token",
		})
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	var req dto.ProjectRequest

	if err = c.BodyParser(&req); err != nil {
		return err
	}

	if err = h.validator.Validate(&req); err != nil {
		return err
	}

	if err = h.projectUC.Edit(c.Context(), projectId, userID, &req); err != nil {
		return err
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(
		fiber.Map{
			"message": "Project success updated",
		},
	)
}

func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	projectIdStr := c.Params("id")

	if projectIdStr == "" {
		return c.Status(
			fiber.StatusNotFound,
		).JSON(fiber.Map{
			"message": "Project ID is empty",
		})
	}

	projectId, err := uuid.Parse(projectIdStr)
	if err != nil {
		return err
	}

	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized: user id not found in token",
		})
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	err = h.projectUC.Remove(c.Context(), projectId, userID)
	if err != nil {
		return err
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(
		fiber.Map{
			"message": "Project success deleted!",
		},
	)
}

func (h *ProjectHandler) DeleteAll(c *fiber.Ctx) error {
	userIdStr, ok := c.Locals("userId").(string)
	if !ok || userIdStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized: user id not found in token",
		})
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	err = h.projectUC.RemoveAll(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(
		fiber.StatusOK,
	).JSON(
		fiber.Map{
			"message": "ALl project success deleted!",
		},
	)
}
