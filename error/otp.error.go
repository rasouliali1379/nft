package error

import "errors"

var (
	ErrInvalidOtpCode = errors.New("invalid OTP code")
)