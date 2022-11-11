package apperrors

import "errors"

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)
