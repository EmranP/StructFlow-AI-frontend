package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/password"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/token"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/verification"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/domain"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/email"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/repository"
	customerrors "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/errors"
	"github.com/google/uuid"
)

type authUseCase struct {
	userRepo   repository.UserRepository
	verifyRepo repository.VerificationRepository

	passwordService password.Service
	tokenService    *token.Service
	emailService    *email.ResendService
}

func New(
	userRepo repository.UserRepository,
	verifyRepo repository.VerificationRepository,
	passwordService password.Service,
	tokenService *token.Service,
	emailService *email.ResendService,
) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		verifyRepo: verifyRepo,

		passwordService: passwordService,
		tokenService:    tokenService,
		emailService:    emailService,
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
	existing, err := u.userRepo.GetByEmail(ctx, email)

	if err == nil && existing != nil {

		if existing.Verify {
			return customerrors.ErrUserAlreadyExists
		}

		verifyCode, err := u.verifyRepo.GetByUserID(ctx, existing.ID)
		if err != nil {
			return err
		}

		if time.Now().UTC().Before(verifyCode.LastSentAt.UTC()) {
			return customerrors.ErrVerificationCooldown
		}

		newCode := verification.GenerateCode()
		err = u.verifyRepo.UpdateCode(
			ctx,
			existing.ID,
			newCode,
			time.Now().UTC().Add(time.Minute),
			time.Now().UTC().Add(15*time.Minute),
		)
		if err != nil {
			return err
		}

		errEmail := u.emailService.SendVerification(
			existing.Email,
			newCode,
		)
		if errEmail != nil {
			fmt.Println(errEmail.Error())
			return errEmail
		}

		return customerrors.ErrVerificationCodeResent
	}

	hash, errHash := u.passwordService.Hash(password)
	if errHash != nil {
		return err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		Verify:       false,
		PasswordHash: hash,
	}

	errNewUser := u.userRepo.Create(ctx, user)
	if errNewUser != nil {
		return err
	}

	verifyCode := &domain.VerificationCode{
		ID:         uuid.New(),
		UserID:     user.ID,
		Code:       verification.GenerateCode(),
		LastSentAt: time.Now().UTC().Add(time.Minute),
		ExpiresAt:  time.Now().UTC().Add(15 * time.Minute),
	}

	errNewVerify := u.verifyRepo.Create(ctx, verifyCode)
	if errNewVerify != nil {
		return errNewVerify
	}

	errEmail := u.emailService.SendVerification(user.Email, verifyCode.Code)
	if errEmail != nil {
		fmt.Println(errEmail.Error())
		return errEmail
	}

	return nil
}

func (u *authUseCase) VerifyEmail(
	ctx context.Context,
	email string,
	code string,
) error {
	user, err := u.userRepo.
		GetByEmail(ctx, email)

	if err != nil {
		return err
	}

	if user == nil {
		return customerrors.ErrUserNotFound
	}

	verifyCode, err := u.verifyRepo.
		GetByUserID(
			ctx,
			user.ID,
		)

	if err != nil {
		return err
	}

	if verifyCode == nil {
		return customerrors.ErrInvalidVerificationCode
	}

	if verifyCode.Code != code {
		return customerrors.ErrInvalidVerificationCode
	}

	if time.Now().After(
		verifyCode.ExpiresAt,
	) {
		return customerrors.ErrVerificationCodeResent
	}

	err = u.userRepo.UpdateVerified(
		ctx,
		user.ID,
	)

	if err != nil {
		return err
	}

	err = u.verifyRepo.Delete(
		ctx,
		verifyCode.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *authUseCase) Me(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	return user.ID.String(), nil
}
