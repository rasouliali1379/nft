package model

import (
	"github.com/google/uuid"
	collection "nft/internal/collection/model"
	nft "nft/internal/nft/model"
	offer "nft/internal/offer/model"
	usermodel "nft/internal/user/model"
	"time"
)

type Sale struct {
	ID            *uuid.UUID
	User          usermodel.User
	Expiration    time.Time
	CanceledBy    *usermodel.User
	CanceledAt    *time.Time
	Collection    *collection.Collection
	Nft           *nft.Nft
	AcceptedOffer *offer.Offer
	MinPrice      float64
	SaleType      Type
	AssetType     AssetType
	Status        Status
}

type Type string

const (
	SaleTypeP2P     Type = "p2p"
	SaleTypeAuction Type = "auction"
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
