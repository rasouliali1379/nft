package apperrors

import "errors"

var (
	ErrNationalIdExists  = errors.New("national id already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrEmailNotFound       = errors.New("email not found")
	ErrPhoneNumberExists = errors.New("phone number already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailDoesntBelongToUser = errors.New("email doesn't belong to user")
)
