package apperrors

import "errors"

var (
	ErrCollectionNotFound      = errors.New("collection doesn't exist")
	ErrCollectionDraftNotFound = errors.New("collection draft not found")
	ErrCollectionIsNotDraft    = errors.New("collection isn't drafted")
	ErrInvalidCollectionId     = errors.New("invalid collection id")
)
