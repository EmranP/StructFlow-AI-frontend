package domain

import "time"

type Template struct {
	ID string

	GenerationID string

	Name string

	Description string

	StructureJSON []byte

	CreatedAt time.Time
	UpdatedAt time.Time
}
