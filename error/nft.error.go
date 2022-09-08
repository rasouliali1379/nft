package apperrors

import "errors"

var (
	ErrNftNotFound              = errors.New("nftasset doesn't exist")
	ErrNftDraftNotFound         = errors.New("nftasset draft not found")
	ErrNftIsNotDraft            = errors.New("nftasset isn't drafted")
	ErrInvalidNftId             = errors.New("invalid nftasset id")
	ErrNftNotSubmittedForReview = errors.New("nftasset is not submitted for review")
)
