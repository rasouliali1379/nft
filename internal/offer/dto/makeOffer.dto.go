package dto

import "github.com/google/uuid"

type MakeOfferRequest struct {
	SaleId uuid.UUID `json:"sale_id" validate:"required"`
	Price  float64   `json:"price" validate:"required"`
}
