package apperrors

import "errors"

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrTokenMalformed = errors.New("malformed token")
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenInvoked = errors.New("token already invoked")
)
