package apperrors

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrNoQueries = errors.New("you need to provide a query")
	ErrInvalidUUID = errors.New("invalid uuid")
)
