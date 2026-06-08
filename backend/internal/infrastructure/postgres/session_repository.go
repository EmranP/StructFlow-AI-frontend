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

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(
	db *pgxpool.Pool,
) *SessionRepository {

	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(
	ctx context.Context,
	token *domain.SessionToken,
) error {

	query := `
		INSERT INTO public.sessions(
			id,
			user_id,
			token_hash,
			expires_at
		)
		VALUES($1,$2,$3,$4)
	`

	_, err := r.db.Exec(
		ctx,
		query,

		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	)

	return err
}

func (r *SessionRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.SessionToken, error) {

	query := `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			created_at
		FROM public.sessions
		WHERE id = $1
	`

	var session domain.SessionToken

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, customerrors.ErrUnauthorized
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) GetByHash(
	ctx context.Context,
	hash string,
) (*domain.SessionToken, error) {
	query := `
		SELECT 
			id,
			user_id,
			token_hash,
			expires_at,
			created_at
		FROM public.sessions
		WHERE token_hash = $1
	`

	var token domain.SessionToken

	err := r.db.QueryRow(ctx, query, hash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, customerrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *SessionRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (uuid.UUID, error) {
	query := `
		SELECT 			
			id
		FROM public.sessions
		WHERE user_id = $1
	`

	var id uuid.UUID

	err := r.db.QueryRow(ctx, query, userID).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, customerrors.ErrSessionNotFound
	}

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *SessionRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	token *domain.SessionToken,
) error {
	query := `
		UPDATE public.sessions SET
		token_hash = $1,
		expires_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, token.TokenHash, token.ExpiresAt, id)

	return err
}

func (r *SessionRepository) DeleteByUserID(
	ctx context.Context,
	userId uuid.UUID,
) error {
	query := `
		DELETE FROM public.sessions
		WHERE user_id = $1
	`

	_, err := r.db.Exec(ctx, query, userId)

	return err
}
