package apperrors

import "errors"

var (
	ErrAppealNotFoundError = errors.New("appeal not found")
	ErrInvalidAppealId     = errors.New("invalid appeal id")
)
