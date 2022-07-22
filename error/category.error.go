package apperrors

import "errors"

var (
	ErrParentCategoryNotFound = errors.New("parent category wasn't found")
	ErrCategoryNotFound = errors.New("category wasn't found")
)
