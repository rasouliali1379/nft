package sale

import "github.com/google/uuid"

type SaleRequest struct {
	AssetId  uuid.UUID `json:"asset_id" validate:"required"`
	MinPrice float64   `json:"min_price" validate:"required"`
	SaleType SaleType  `json:"sale_type" validate:"required"`
}

type SaleResponse struct {
	SaleId     uuid.UUID `json:"sale_id"`
	Expiration int64     `json:"expiration"`
}
