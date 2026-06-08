package repository

import (
	"context"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(
		ctx context.Context,
		token *domain.SessionToken,
	) error

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.SessionToken, error)

	GetByHash(
		ctx context.Context,
		hash string,
	) (*domain.SessionToken, error)

	GetByUserID(
		ctx context.Context,
		userId uuid.UUID,
	) (uuid.UUID, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		token *domain.SessionToken,
	) error

	DeleteByUserID(
		ctx context.Context,
		userId uuid.UUID,
	) error
}
