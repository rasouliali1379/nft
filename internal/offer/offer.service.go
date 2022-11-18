package offer

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
	"nft/internal/offer/model"
)

type OfferService struct {
	offerRepository contract.IOfferRepository
	saleRepository  contract.ISaleRepository
}

type OfferServiceParams struct {
	fx.In
	OfferRepository contract.IOfferRepository
	SaleRepository  contract.ISaleRepository
}

func NewOfferService(params OfferServiceParams) contract.IOfferService {
	return &OfferService{
		offerRepository: params.OfferRepository,
		saleRepository:  params.SaleRepository,
	}
}

func (o OfferService) MakeOfferToSale(c context.Context, m model.Offer) error {
	span, c := jtrace.T().SpanFromContext(c, "OfferService[MakeOfferToSale]")
	defer span.Finish()

	sale, err := o.saleRepository.Get(c, persist.D{"id": m.SaleId})
	if err != nil {
		return err
	}

	if sale.MinPrice > m.Price {
		return apperrors.ErrOfferLowerMinPrice
	}

	if sale.User.ID == m.User.ID {
		return apperrors.ErrOfferYourSale
	}

	_, err = o.offerRepository.Add(c, m)
	if err != nil {
		return err
	}

	return nil
}

func (o OfferService) CancelOffer(c context.Context, m model.Offer) error {
	span, c := jtrace.T().SpanFromContext(c, "OfferService[CancelOffer]")
	defer span.Finish()

	offerModel, err := o.offerRepository.Get(c, persist.D{"id": *m.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrOfferNotFound
		}
		return err
	}

	if offerModel.User.ID != m.User.ID {
		return apperrors.ErrOfferNotFound
	}

	return o.offerRepository.Delete(c, model.Offer{ID: m.ID})
}

func (o OfferService) AcceptOffer(c context.Context, m model.Offer) error {
	span, c := jtrace.T().SpanFromContext(c, "OfferService[AcceptOffer]")
	defer span.Finish()

	offerModel, err := o.offerRepository.Get(c, persist.D{"id": *m.ID})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrOfferNotFound
		}
		return err
	}

	sale, err := o.saleRepository.Get(c, persist.D{"id": offerModel.SaleId})
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return apperrors.ErrSaleNotFound
		}
		return err
	}

	if m.User.ID != sale.User.ID {
		return apperrors.ErrOfferNotFound
	}

	_, err = o.offerRepository.Update(c, model.Offer{ID: m.ID, Accepted: true})
	if err != nil {
		return err
	}

	return nil
}

func (o OfferService) GetAllOffers(c context.Context, m model.Offer) ([]model.Offer, error) {
	span, c := jtrace.T().SpanFromContext(c, "OfferService[GetAllOffers]")
	defer span.Finish()
	return o.offerRepository.GetAll(c, persist.D{"sale_id": m.SaleId})
}
