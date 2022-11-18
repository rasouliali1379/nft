package sale

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
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

	createdSale, err := s.db.Create(c, &saleEntity)
	if err != nil {
		return model.Sale{}, err
	}

	return mapSaleEntityToModel(*createdSale.(*entity.Sale)), nil
}

func (s SaleRepository) Get(c context.Context, conditions persist.D) (model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleRepository[Get]")
	defer span.Finish()

	sale, err := s.db.Get(c, &entity.Sale{}, conditions)
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return model.Sale{}, apperrors.ErrSaleNotFound
		}
		return model.Sale{}, err
	}

	return mapSaleEntityToModel(*sale.(*entity.Sale)), nil
}

func (s SaleRepository) Cancel(c context.Context, m model.Sale) error {
	span, c := jtrace.T().SpanFromContext(c, "SaleRepository[Cancel]")
	defer span.Finish()

	data := persist.D{"canceled_at": time.Now(), "canceled_by": m.User.ID}
	if _, err := s.db.Update(c, &entity.Sale{ID: *m.ID}, data); err != nil {
		return err
	}
	return nil
}

func (s SaleRepository) GetAll(c context.Context, conditions persist.D) ([]model.Sale, error) {
	span, c := jtrace.T().SpanFromContext(c, "SaleRepository[GetAll]")
	defer span.Finish()

	saleList, err := s.db.GetAll(c, &[]entity.Sale{}, conditions)
	if err != nil {
		return nil, err
	}

	return createModelSaleListFromEntity(*saleList.(*[]entity.Sale)), nil
}
