package apperrors

import "errors"

var (
	ErrOfferLowerMinPrice = errors.New("offer price is should be higher than min price")
	ErrOfferYourSale      = errors.New("you can't make offer on your sale")
	ErrOfferNotFound      = errors.New("offer not found")
)
