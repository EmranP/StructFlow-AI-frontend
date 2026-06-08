package domain

import (
	"time"

	"github.com/google/uuid"
)

type VerificationCode struct {
	ID uuid.UUID

	UserID uuid.UUID

	Code string

	LastSentAt time.Time
	ExpiresAt  time.Time

	CreatedAt time.Time
}
