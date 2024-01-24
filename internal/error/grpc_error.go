package grpcerror

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNoToken      = errors.New("authorization token was not provided")
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenClaims  = errors.New("failed to get token claims")
	ErrForbidden    = errors.New("forbidden")
)
