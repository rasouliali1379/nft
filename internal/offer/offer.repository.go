package offer

import (
	"go.uber.org/fx"
	"nft/contract"
)

type OfferRepository struct {
	db contract.IPersist
}

type OfferRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewOfferRepository(params OfferRepositoryParams) contract.IOfferRepository {
	return &OfferRepository{
		db: params.DB,
	}
}
