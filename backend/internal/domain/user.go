package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID

	Email        string
	PasswordHash string
	Verify       bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
