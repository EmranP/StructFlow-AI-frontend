package repository

import (
	"context"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(
		ctx context.Context,
		project *domain.Project,
	) error

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Project, error)

	GetAllByUserID(
		ctx context.Context,
		userID uuid.UUID,
		limit,
		offset int,
	) ([]domain.Project, error)

	GetCount(ctx context.Context, userID uuid.UUID) (int64, error)

	Update(
		ctx context.Context,
		id,
		userID uuid.UUID,
		data *domain.Project,
	) error

	Delete(
		ctx context.Context,
		id,
		userID uuid.UUID,
	) error

	DeleteAll(
		ctx context.Context,
		userID uuid.UUID,
	) error
}
