package usecase

import (
	"context"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/token"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(
		ctx context.Context,
		email string,
		password string,
	) error

	Login(
		ctx context.Context,
		email string,
		password string,
	) (*token.GenerationTokens, error)

	VerifyEmail(
		ctx context.Context,
		email string,
		code string,
	) (*token.GenerationTokens, error)

	Me(
		ctx context.Context,
		id uuid.UUID,
	) (string, error)
}

type AuthSessionUseCase interface {
	Generate(userID uuid.UUID, sessionID uuid.UUID) (*token.GenerationTokens, *domain.SessionToken, error)

	Create(ctx context.Context, token *domain.SessionToken) error
	FindByUserID(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	Update(ctx context.Context, userId uuid.UUID, token *domain.SessionToken) error
	CreateOrUpdate(
		ctx context.Context,
		userID uuid.UUID,
	) (*token.GenerationTokens, error)

	Refresh(
		ctx context.Context,
		refreshToken string,
	) (*token.GenerationTokens, error)
	Clear(
		ctx context.Context,
		userId uuid.UUID,
	) error
}
