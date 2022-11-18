package sale

import (
	collection "nft/internal/collection/dto"
	nft "nft/internal/nft/dto"
)

type SaleList struct {
	Sales []Sale `json:"sales"`
}

type Sale struct {
	ID         string                 `json:"id"`
	Expiration int64                  `json:"expiration"`
	Collection *collection.Collection `json:"collection,omitempty"`
	Nft        *nft.Nft               `json:"nft,omitempty"`
	MinPrice   float64                `json:"min_price"`
	SaleType   SaleType               `json:"sale_type"`
	AssetType  AssetType              `json:"asset_type"`
	Status     Status                 `json:"status"`
}

type SaleType string

const (
	SaleTypeP2P     SaleType = "p2p"
	SaleTypeAuction SaleType = "auction"
)

type AssetType string

const (
	AssetTypeNft        AssetType = "nft"
	AssetTypeCollection AssetType = "collection"
)

type Status string

const (
	SaleStatusSold       Status = "sold"
	SaleStatusInProgress Status = "in_progress"
	SaleStatusCanceled   Status = "canceled"
	SaleStatusExpired    Status = "expired"
)
