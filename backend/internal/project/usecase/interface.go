package usecase

import (
	"context"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/project/dto"
	"github.com/google/uuid"
)

type ProjectUseCase interface {
	Add(
		ctx context.Context,
		userID uuid.UUID,
		p *dto.ProjectRequest,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Project, error)

	FindByUserID(
		ctx context.Context,
		userID uuid.UUID,
		page,
		limit int,
	) ([]dto.ProjectResponse, int64, error)

	Edit(
		ctx context.Context,
		id,
		userID uuid.UUID,
		data *dto.ProjectRequest,
	) error

	Remove(
		ctx context.Context,
		id,
		userID uuid.UUID,
	) error

	RemoveAll(
		ctx context.Context,
		userID uuid.UUID,
	) error
}
