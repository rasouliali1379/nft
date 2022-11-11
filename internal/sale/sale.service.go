package sale

import (
	"context"
	"go.uber.org/fx"
	"nft/contract"
	"nft/infra/jtrace"
	nftmodel "nft/internal/nft/model"
	"nft/internal/sale/model"
	usermodel "nft/internal/user/model"
)

type SaleService struct {
	saleRepository contract.ISaleRepository
	nftService     contract.INftService
}

type SaleServiceParams struct {
	fx.In
	SaleRepository contract.ISaleRepository
	NftService     contract.INftService
}

func NewSaleService(params SaleServiceParams) contract.ISaleService {
	return SaleService{
		saleRepository: params.SaleRepository,
		nftService:     params.NftService,
	}
}

func (s SaleService) CreateNftSale(c context.Context, m model.Sale) (model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleService[CreateNftSale]")
	defer span.Finish()

	nft, err := s.nftService.GetOwnedNft(c, nftmodel.Nft{ID: &m.AssetId, CurrentOwner: &m.User})
	if err != nil {
		return model.Sale{}, err
	}

	m.User = usermodel.User{ID: nft.CurrentOwner.ID}
	m.SaleType = model.SaleTypeP2P
	m.AssetType = model.AssetTypeNft

	return s.saleRepository.Create(c, m)
}

func (s SaleService) CreateCollectionSale(c context.Context, m model.Sale) (model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleService[CreateCollectionSale]")
	defer span.Finish()

	nft, err := s.nftService.GetOwnedNft(c, nftmodel.Nft{ID: &m.AssetId, CurrentOwner: &m.User})
	if err != nil {
		return model.Sale{}, err
	}

	m.User = usermodel.User{ID: nft.CurrentOwner.ID}
	m.SaleType = model.SaleTypeP2P
	m.AssetType = model.AssetTypeCollection

	return s.saleRepository.Create(c, m)
}
