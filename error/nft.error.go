package apperrors

import "errors"

var (
	ErrNftNotFound              = errors.New("nft doesn't exist")
	ErrNftDraftNotFound         = errors.New("nft draft not found")
	ErrNftIsNotDraft            = errors.New("nft isn't drafted")
	ErrInvalidNftId             = errors.New("invalid nft id")
	ErrNftNotSubmittedForReview = errors.New("nft is not submitted for review")
)
