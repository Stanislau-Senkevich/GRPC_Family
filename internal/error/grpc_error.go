package grpcerror

import "errors"

var (
	ErrInviteExist     = errors.New("user already invited")
	ErrInviteNotFound  = errors.New("invite not found")
	ErrUserInFamily    = errors.New("user already in family")
	ErrUserNotInFamily = errors.New("user not in family")
	ErrInternalError   = errors.New("internal error")
	ErrUserNotFound    = errors.New("user not found")
	ErrFamilyNotFound  = errors.New("family not found")
	ErrNoToken         = errors.New("authorization token was not provided")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenClaims     = errors.New("failed to get token claims")
	ErrForbidden       = errors.New("forbidden")
)
