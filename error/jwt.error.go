package error

import "errors"

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrTokenMalformed = errors.New("malformed token")
	ErrTokenInvoked = errors.New("token invoked")
)
