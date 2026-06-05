package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/dto"
)

type GenerationUseCase interface {
	Add(
		ctx context.Context,
		projectID uuid.UUID,
	) (*domain.Generation, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Generation, error)

	FindByProjectID(
		ctx context.Context,
		projectId uuid.UUID,
		limit,
		page int,
	) ([]dto.GenerationAllResponse, int64, error)

	FindTemplates(
		ctx context.Context,
		genId uuid.UUID,
	) ([]dto.GenerationTempResponse, error)

	Download(
		ctx context.Context,
		generationID uuid.UUID,
	) ([]byte, error)
}
