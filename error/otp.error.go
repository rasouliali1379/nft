package apperrors

import "errors"

var (
	ErrInvalidOtpCode = errors.New("invalid OTP code")
)