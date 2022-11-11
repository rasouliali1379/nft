package sale

import "github.com/google/uuid"

type SaleRequest struct {
	NftId    uuid.UUID `json:"nft_id" validate:"required"`
	MinPrice float64   `json:"min_price" validate:"required"`
}

type SaleResponse struct {
	SaleId     uuid.UUID `json:"sale_id"`
	Expiration int64     `json:"expiration"`
}
