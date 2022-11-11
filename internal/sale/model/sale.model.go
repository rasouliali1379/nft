package model

import (
	"github.com/google/uuid"
	usermodel "nft/internal/user/model"
	"time"
)

type Sale struct {
	ID         *uuid.UUID
	User       usermodel.User
	Expiration time.Time
	CanceledBy *usermodel.User
	CanceledAt *time.Time
	AssetId    uuid.UUID
	MinPrice   float64
	SaleType   SaleType
	AssetType  AssetType
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
