package repository

import (
	"context"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user *domain.User,
	) error

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.User, error)

	GetByEmail(
		ctx context.Context,
		email string,
	) (*domain.User, error)

	UpdateVerified(
		ctx context.Context,
		id uuid.UUID,
	) error

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
