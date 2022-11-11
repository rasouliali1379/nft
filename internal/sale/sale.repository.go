package sale

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/contract"
	"nft/infra/jtrace"
	entity "nft/internal/sale/entity"
	"nft/internal/sale/model"
	"time"
)

type SaleRepository struct {
	db contract.IPersist
}

type SaleRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewSaleRepository(params SaleRepositoryParams) contract.ISaleRepository {
	return &SaleRepository{
		db: params.DB,
	}
}

func (s SaleRepository) Create(c context.Context, m model.Sale) (model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleRepository[Create]")
	defer span.Finish()

	saleEntity := mapSaleModelToEntity(m)
	saleEntity.ID = uuid.New()
	saleEntity.Expiration = time.Now().Add(time.Hour * 168)

	createdNft, err := s.db.Create(c, &saleEntity)
	if err != nil {
		return model.Sale{}, err
	}

	return mapSaleEntityToModel(*createdNft.(*entity.Sale)), nil
}
