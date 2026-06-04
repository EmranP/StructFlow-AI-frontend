package usecase

import (
	"context"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/password"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/token"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/repository"
	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
	"github.com/google/uuid"
)

type authUseCase struct {
	userRepo repository.UserRepository

	passwordService password.Service
	tokenService    *token.Service
}

func New(
	userRepo repository.UserRepository,
	passwordService password.Service,
	tokenService *token.Service,
) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,

		passwordService: passwordService,
		tokenService:    tokenService,
	}
}

func (u *authUseCase) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := u.userRepo.
		GetByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", customerrors.ErrInvalidCredentials
	}

	if !u.passwordService.Verify(
		password,
		user.PasswordHash,
	) {
		return "", customerrors.ErrInvalidCredentials
	}

	token, err := u.tokenService.
		GenerateAccessToken(
			user.ID.String(),
		)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUseCase) Register(
	ctx context.Context,
	email string,
	password string,
) error {

	existing, err := u.userRepo.
		GetByEmail(ctx, email)

	if err != nil {
		return err
	}

	if existing != nil {
		return customerrors.ErrUserAlreadyExists
	}

	hash, err := u.passwordService.
		Hash(password)

	if err != nil {
		return err
	}

	user := &domain.User{
		ID: uuid.New(),

		Email: email,

		PasswordHash: hash,
	}

	return u.userRepo.Create(
		ctx,
		user,
	)
}
