package sale

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	persist "nft/infra/persist/type"
	collection "nft/internal/collection/model"
	nft "nft/internal/nft/model"
	"nft/internal/sale/model"
	usermodel "nft/internal/user/model"
)

type SaleService struct {
	saleRepository    contract.ISaleRepository
	nftService        contract.INftService
	collectionService contract.ICollectionService
	offerRepository   contract.IOfferRepository
}

type SaleServiceParams struct {
	fx.In
	SaleRepository    contract.ISaleRepository
	NftService        contract.INftService
	CollectionService contract.ICollectionService
	OfferRepository   contract.IOfferRepository
}

func NewSaleService(params SaleServiceParams) contract.ISaleService {
	return &SaleService{
		saleRepository:    params.SaleRepository,
		nftService:        params.NftService,
		collectionService: params.CollectionService,
		offerRepository:   params.OfferRepository,
	}
}

func (s SaleService) CreateNftSale(c context.Context, m model.Sale) (model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleService[CreateNftSale]")
	defer span.Finish()

	ownedNft, err := s.nftService.GetOwnedNft(c, nft.Nft{ID: m.Nft.ID, CurrentOwner: &m.User})
	if err != nil {
		return model.Sale{}, err
	}

	m.User = usermodel.User{ID: ownedNft.CurrentOwner.ID}
	m.AssetType = model.AssetTypeNft

	return s.saleRepository.Create(c, m)
}

func (s SaleService) CreateCollectionSale(c context.Context, m model.Sale) (model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleService[CreateCollectionSale]")
	defer span.Finish()

	nftModel, err := s.nftService.GetOwnedNft(c, nft.Nft{ID: m.Nft.ID, CurrentOwner: &m.User})
	if err != nil {
		return model.Sale{}, err
	}

	m.User = usermodel.User{ID: nftModel.CurrentOwner.ID}
	m.AssetType = model.AssetTypeCollection

	return s.saleRepository.Create(c, m)
}

func (s SaleService) CancelSale(c context.Context, m model.Sale) error {
	span, c := jtrace.T().SpanFromContext(c, "SaleService[CancelSale]")
	defer span.Finish()

	_, err := s.saleRepository.Get(c, persist.D{"id": *m.ID, "user_id": m.User.ID})
	if err != nil {
		return err
	}

	return s.saleRepository.Cancel(c, m)
}

func (s SaleService) GetSalesList(c context.Context, userId uuid.UUID) ([]model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleService[GetSalesList]")
	defer span.Finish()

	sales, err := s.saleRepository.GetAll(c, persist.D{"user_id": userId})
	if err != nil {
		return nil, err
	}

	for i := range sales {
		switch sales[i].AssetType {
		case model.AssetTypeNft:
			nftModel, err := s.nftService.GetNft(c, nft.Nft{ID: sales[i].Nft.ID})
			if err != nil {
				return nil, err
			}
			sales[i].Nft = &nftModel
		case model.AssetTypeCollection:
			collectionModel, err := s.collectionService.GetCollection(c, collection.Collection{ID: sales[i].Collection.ID})
			if err != nil {
				return nil, err
			}
			sales[i].Collection = &collectionModel
		}

		offer, err := s.offerRepository.Get(c, persist.D{"sale_id": *sales[i].ID, "accepted": true})
		if err != nil {
			if errors.Is(err, apperrors.ErrOfferNotFound) {
				sales[i].Status = model.SaleStatusInProgress
			} else {
				return nil, err
			}
		}

		if offer.Accepted {
			sales[i].Status = model.SaleStatusSold
			sales[i].AcceptedOffer = &offer
		}
	}

	return sales, nil
}

func (s SaleService) GetSale(c context.Context, m model.Sale) (model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleService[GetSale]")
	defer span.Finish()

	sale, err := s.saleRepository.Get(c, persist.D{"id": *m.ID, "user_id": m.User.ID})
	if err != nil {
		return model.Sale{}, err
	}

	switch sale.AssetType {
	case model.AssetTypeNft:
		nftModel, err := s.nftService.GetNft(c, nft.Nft{ID: sale.Nft.ID})
		if err != nil {
			return model.Sale{}, err
		}
		sale.Nft = &nftModel
	case model.AssetTypeCollection:
		collectionModel, err := s.collectionService.GetCollection(c, collection.Collection{ID: sale.Collection.ID})
		if err != nil {
			return model.Sale{}, err
		}
		sale.Collection = &collectionModel
	}

	offer, err := s.offerRepository.Get(c, persist.D{"sale_id": *sale.ID, "deleted_at": nil, "accepted": true})
	if err != nil {
		log.Println(err)
		if errors.Is(err, apperrors.ErrOfferNotFound) {
			sale.Status = model.SaleStatusInProgress
		} else {
			return model.Sale{}, err
		}
	}

	if offer.Accepted {
		sale.Status = model.SaleStatusSold
		sale.AcceptedOffer = &offer
	}

	return sale, nil
}
