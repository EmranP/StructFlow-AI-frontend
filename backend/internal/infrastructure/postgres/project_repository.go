package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(
	db *pgxpool.Pool,
) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

const getProjectSQL = "id, title, project_type, stack, architecture, features, additional_info, created_at, updated_at"

func (r *ProjectRepository) Create(
	ctx context.Context,
	project *domain.Project,
) (uuid.UUID, error) {

	query := `
		INSERT INTO public.projects(
			id,
			user_id,
			title,
			project_type,
			stack,
			architecture,
			features,
			additional_info
		)
		VALUES(
			$1,$2,$3,$4,$5,$6,$7,$8
		)
		RETURNING id;
	`

	var projectId uuid.UUID

	err := r.db.QueryRow(
		ctx,
		query,

		project.ID,
		project.UserID,
		project.Title,
		project.ProjectType,
		project.Stack,
		project.Architecture,
		project.Features,
		project.AdditionalInfo,
	).Scan(
		&projectId,
	)

	if err != nil {
		return uuid.Nil, err
	}

	return projectId, err
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM public.projects
		WHERE id = $1;
	`, getProjectSQL)

	var project domain.Project

	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Title,
		&project.ProjectType,
		&project.Stack,
		&project.Architecture,
		&project.Features,
		&project.AdditionalInfo,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, customerrors.ErrProjectNotFound
	}

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepository) GetAllByUserID(ctx context.Context,
	userID uuid.UUID,
	limit,
	offset int,
) ([]domain.Project, error) {
	query := fmt.Sprintf(`
			SELECT %s FROM public.projects
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT $2
			OFFSET $3;
	`, getProjectSQL)

	projects := make([]domain.Project, 0)

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Project

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.ProjectType,
			&p.Stack,
			&p.Architecture,
			&p.Features,
			&p.AdditionalInfo,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) GetCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM public.projects WHERE user_id = $1`

	var totalCount int64

	err := r.db.QueryRow(ctx, query, userID).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r *ProjectRepository) Update(
	ctx context.Context,
	id,
	userID uuid.UUID,
	data *domain.Project,
) error {
	query := `
		UPDATE public.projects SET
		title = COALESCE($1, title),
		project_type = COALESCE($2, project_type),
		stack = COALESCE($3, stack),
		architecture = COALESCE($4, architecture),
		features = COALESCE($5, features),
		additional_info = COALESCE($6, additional_info)
		WHERE id = $7 AND user_id = $8
	`

	_, err := r.db.Exec(
		ctx,
		query,
		&data.Title,
		&data.ProjectType,
		&data.Stack,
		&data.Architecture,
		&data.Features,
		&data.AdditionalInfo,
		id,
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) Delete(
	ctx context.Context,
	id,
	userID uuid.UUID,
) error {
	query := `
		DELETE FROM public.projects 
		WHERE id = $1 AND user_id = $2;
	`

	_, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) DeleteAll(
	ctx context.Context,
	userID uuid.UUID,
) error {
	query := `DELETE FROM public.projects WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
