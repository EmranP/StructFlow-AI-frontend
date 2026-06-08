package repository

import (
	"context"
	"time"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/google/uuid"
)

type VerificationRepository interface {
	Create(
		ctx context.Context,
		code *domain.VerificationCode,
	) error

	GetByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) (*domain.VerificationCode, error)

	UpdateCode(
		ctx context.Context,
		userID uuid.UUID,
		code string,
		lastSentAt time.Time,
		expiresAt time.Time,
	) error

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
