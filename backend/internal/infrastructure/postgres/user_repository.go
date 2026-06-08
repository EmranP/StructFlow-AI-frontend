package postgres

import (
	"context"
	"errors"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(
	db *pgxpool.Pool,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {

	query := `
		INSERT INTO public.users(
			id,
			email,
			password_hash
		)
		VALUES($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,

		user.ID,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {

	query := `
		SELECT
			id,
			email,
			password_hash,
			is_verified,
			created_at,
			updated_at
		FROM public.users
		WHERE email = $1
	`

	var user domain.User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Verify,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, customerrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {

	query := `
		SELECT
			id,
			email,
			password_hash,
			is_verified,
			created_at,
			updated_at
		FROM public.users
		WHERE id = $1
	`

	var user domain.User

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Verify,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateVerified(
	ctx context.Context,
	id uuid.UUID,
) error {

	query := `
		UPDATE public.users
		SET is_verified = true
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}

func (r *UserRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
		DELETE FROM public.users
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)

	return err
}
