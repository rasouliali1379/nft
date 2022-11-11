package entity

import (
	"github.com/google/uuid"
	"time"
)

type Sale struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt *time.Time

	UserId     uuid.UUID
	Expiration time.Time
	CanceledBy *uuid.UUID
	CanceledAt *time.Time
	SaleType   SaleType
	AssetType  AssetType
	AssetId    uuid.UUID
	MinPrice   float64
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
