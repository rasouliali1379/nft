package offer

import (
	"go.uber.org/fx"
	"nft/contract"
)

type OfferService struct {
	offerRepository contract.IOfferRepository
}

type OfferServiceParams struct {
	fx.In
	OfferRepository contract.IOfferRepository
}

func NewOfferService(params OfferServiceParams) contract.IOfferService {
	return OfferService{
		offerRepository: params.OfferRepository,
	}
}
