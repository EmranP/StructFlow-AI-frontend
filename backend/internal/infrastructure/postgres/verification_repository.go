package postgres

import (
	"context"
	"time"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationRepository struct {
	db *pgxpool.Pool
}

func NewVerificationRepository(
	db *pgxpool.Pool,
) *VerificationRepository {
	return &VerificationRepository{
		db: db,
	}
}

func (r *VerificationRepository) Create(
	ctx context.Context,
	code *domain.VerificationCode,
) error {

	query := `
		INSERT INTO public.verification_codes(
			id,
			user_id,
			code,
			last_sent_at,
			expires_at
		)
		VALUES($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		code.ID,
		code.UserID,
		code.Code,

		code.LastSentAt,
		code.ExpiresAt,
	)

	return err
}

func (r *VerificationRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*domain.VerificationCode, error) {

	query := `
		SELECT
			id,
			user_id,
			code,
			expires_at,
			last_sent_at,
			created_at
		FROM public.verification_codes
		WHERE user_id = $1
		LIMIT 1
	`

	var code domain.VerificationCode

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&code.ID,
		&code.UserID,
		&code.Code,
		&code.ExpiresAt,
		&code.LastSentAt,
		&code.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &code, nil
}

func (r *VerificationRepository) UpdateCode(
	ctx context.Context,
	userID uuid.UUID,
	code string,
	lastSentAt time.Time,
	expiresAt time.Time,
) error {

	query := `
		UPDATE public.verification_codes
		SET
			code = $1,
			last_sent_at = $2,
			expires_at = $3
		WHERE user_id = $4
	`

	_, err := r.db.Exec(
		ctx,
		query,
		code,
		lastSentAt,
		expiresAt,
		userID,
	)

	return err
}

func (r *VerificationRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	query := `
		DELETE FROM public.verification_codes
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}
