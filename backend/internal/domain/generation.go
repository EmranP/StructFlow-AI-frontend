package domain

import "time"

type Generation struct {
	ID string

	ProjectID string

	Status string

	CreatedAt time.Time
	UpdatedAt time.Time
}
