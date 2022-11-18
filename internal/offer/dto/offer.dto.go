package dto

import user "nft/internal/user/dto"

type OfferList struct {
	Offers []Offer `json:"offers"`
}

type Offer struct {
	ID    string    `json:"id"`
	Price float64   `json:"price"`
	User  user.User `json:"user"`
}
