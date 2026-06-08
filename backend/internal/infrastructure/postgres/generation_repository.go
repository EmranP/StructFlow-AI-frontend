package postgres

import (
	"context"
	"errors"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GenerationRepository struct {
	db *pgxpool.Pool
}

func NewGenerationRepository(
	db *pgxpool.Pool,
) *GenerationRepository {
	return &GenerationRepository{
		db: db,
	}
}

func (r *GenerationRepository) Create(
	ctx context.Context,
	generation *domain.Generation,
) error {

	query := `
		INSERT INTO public.generations(
			id,
			project_id,
			status,
			error_message
		)
		VALUES(
			$1,
			$2,
			$3,
			$4
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		generation.ID,
		generation.ProjectID,
		generation.Status,
		generation.ErrorMessage,
	)

	return err
}

func (r *GenerationRepository) GetAllByProjectID(
	ctx context.Context,
	projectID uuid.UUID,
	limit,
	offset int,
) ([]domain.Generation, error) {
	query := `
		SELECT
			id,
			project_id,
			status,
			error_message,
			created_at,
			updated_at
		FROM public.generations
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3;
		`
	gens := make([]domain.Generation, 0)

	rows, err := r.db.Query(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g domain.Generation

		err := rows.Scan(
			&g.ID,
			&g.ProjectID,
			&g.Status,
			&g.ErrorMessage,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {

			return nil, err
		}

		gens = append(gens, g)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return gens, nil
}

func (r *GenerationRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Generation, error) {

	query := `
		SELECT
			id,
			project_id,
			status,
			error_message,
			created_at,
			updated_at
		FROM public.generations
		WHERE id = $1
	`

	var generation domain.Generation

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&generation.ID,
		&generation.ProjectID,
		&generation.Status,
		&generation.ErrorMessage,
		&generation.CreatedAt,
		&generation.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {

		return nil, err
	}

	return &generation, nil
}

func (r *GenerationRepository) GetCount(ctx context.Context, projectID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM public.generations WHERE project_id = $1`

	var totalCount int64

	err := r.db.QueryRow(ctx, query, projectID).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r *GenerationRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
	errorMessage *string,
) error {

	query := `
		UPDATE public.generations
		SET
			status = $2,
			error_message = $3,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
		status,
		errorMessage,
	)

	return err
}
