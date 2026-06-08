package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/password"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/token"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/repository"
	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
	"github.com/google/uuid"
)

type authSessionUseCase struct {
	sessionRepo repository.SessionRepository
	userRepo    repository.UserRepository

	passwordService password.Service
	tokenService    *token.Service
}

func NewSession(
	sessionRepo repository.SessionRepository,
	userRepo repository.UserRepository,
	passwordService password.Service,
	tokenService *token.Service,
) AuthSessionUseCase {
	return &authSessionUseCase{
		sessionRepo:     sessionRepo,
		userRepo:        userRepo,
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}

func (u *authSessionUseCase) Generate(
	sessionID uuid.UUID,
	userID uuid.UUID,
) (*token.GenerationTokens, *domain.SessionToken, error) {

	accessToken, err := u.tokenService.GenerateAccessToken(userID)
	if err != nil {
		return nil, nil, err
	}

	secret, err := token.GenerateRefreshSecret()
	if err != nil {
		return nil, nil, err
	}

	hash, err := u.passwordService.Hash(secret)
	if err != nil {
		return nil, nil, err
	}

	refreshToken := sessionID.String() + "." + secret

	session := &domain.SessionToken{
		ID:        sessionID,
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	return &token.GenerationTokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, session, nil
}

func (u *authSessionUseCase) Create(
	ctx context.Context,
	tokenData *domain.SessionToken,
) error {

	return u.sessionRepo.Create(
		ctx,
		tokenData,
	)
}

func (u *authSessionUseCase) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (uuid.UUID, error) {
	id, err := u.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (u *authSessionUseCase) Update(
	ctx context.Context,
	id uuid.UUID,
	token *domain.SessionToken,
) error {
	err := u.sessionRepo.Update(ctx, id, token)
	if err != nil {
		return err
	}

	return nil
}

func (u *authSessionUseCase) CreateOrUpdate(
	ctx context.Context,
	userID uuid.UUID,
) (*token.GenerationTokens, error) {
	sessionID, err := u.FindByUserID(
		ctx,
		userID,
	)

	create := false

	if errors.Is(err, customerrors.ErrSessionNotFound) {
		sessionID = uuid.New()

		create = true
	} else if err != nil {
		return nil, err
	}

	tokens, session, err := u.Generate(sessionID, userID)
	if err != nil {

		return nil, err
	}

	if create {
		err = u.sessionRepo.Create(ctx, session)
	} else {
		err = u.sessionRepo.Update(
			ctx,
			sessionID,
			session,
		)
	}

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *authSessionUseCase) Refresh(
	ctx context.Context,
	refreshToken string,
) (*token.GenerationTokens, error) {

	sessionID, secret, err := token.ParseRefreshToken(
		refreshToken,
	)
	if err != nil {
		return nil, customerrors.ErrUnauthorized
	}

	session, err := u.sessionRepo.GetByID(
		ctx,
		sessionID,
	)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, customerrors.ErrUnauthorized
	}

	if !u.passwordService.Verify(
		secret,
		session.TokenHash,
	) {
		return nil, customerrors.ErrUnauthorized
	}

	return u.rotateSession(
		ctx,
		session,
	)
}

func (u *authSessionUseCase) rotateSession(
	ctx context.Context,
	session *domain.SessionToken,
) (*token.GenerationTokens, error) {

	tokens, newSession, err := u.Generate(
		session.ID,
		session.UserID,
	)

	if err != nil {
		return nil, err
	}

	err = u.sessionRepo.Update(
		ctx,
		session.ID,
		&domain.SessionToken{
			TokenHash: newSession.TokenHash,
			ExpiresAt: newSession.ExpiresAt,
		},
	)

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *authSessionUseCase) Clear(
	ctx context.Context,
	userId uuid.UUID,
) error {

	return u.sessionRepo.
		DeleteByUserID(
			ctx,
			userId,
		)
}
