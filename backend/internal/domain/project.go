package domain

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Title          string
	ProjectType    string
	Stack          string
	Architecture   string
	Features       string
	AdditionalInfo string
	Prompt         string

	CreatedAt time.Time
	UpdatedAt time.Time
}
