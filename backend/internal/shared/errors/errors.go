package error

import "errors"

var (
	ErrUserAlreadyExists = errors.New(
		"user already exists",
	)

	ErrUserNotFound    = errors.New("user not found")
	ErrProjectNotFound = errors.New("project not found")

	ErrInvalidCredentials = errors.New(
		"invalid credentials",
	)

	ErrUnauthorized = errors.New(
		"unauthorized",
	)
)
