package apperrors

import "errors"

var (
	ErrAppealNotFound  = errors.New("appeal not found")
	ErrInvalidAppealId = errors.New("invalid appeal id")
)
