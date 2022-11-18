package offer

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
	entity "nft/internal/offer/entity"
	"nft/internal/offer/model"
	"time"
)

type OfferRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewOfferRepository(params OfferRepositoryParams) contract.IOfferRepository {
	return &OfferRepository{
		db: params.DB,
	}
}

type OfferRepository struct {
	db contract.IPersist
}

func (o OfferRepository) Add(c context.Context, m model.Offer) (model.Offer, error) {
	span, c := jtrace.T().SpanFromContext(c, "OfferRepository[Add]")
	defer span.Finish()

	offerEntity := mapOfferModelToEntity(m)
	offerEntity.ID = uuid.New()

	createdOffer, err := o.db.Create(c, &offerEntity)
	if err != nil {
		return model.Offer{}, err
	}

	return mapOfferEntityToModel(*createdOffer.(*entity.Offer)), nil
}

func (o OfferRepository) Get(c context.Context, conditions persist.D) (model.Offer, error) {
	span, c := jtrace.T().SpanFromContext(c, "OfferRepository[Get]")
	defer span.Finish()

	offer, err := o.db.Get(c, &entity.Offer{}, conditions)
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return model.Offer{}, apperrors.ErrOfferNotFound
		}
		return model.Offer{}, err
	}

	return mapOfferEntityToModel(*offer.(*entity.Offer)), nil
}

func (o OfferRepository) GetAll(c context.Context, conditions persist.D) ([]model.Offer, error) {
	span, c := jtrace.T().SpanFromContext(c, "OfferRepository[GetAll]")
	defer span.Finish()

	kycList, err := o.db.GetAll(c, &[]entity.Offer{}, conditions)
	if err != nil {
		return nil, err
	}

	return createModelOfferList(*kycList.(*[]entity.Offer)), nil
}

func createModelOfferList(list []entity.Offer) []model.Offer {
	kycList := make([]model.Offer, len(list))

	for i := range list {
		kycList[i] = mapOfferEntityToModel(list[i])
	}

	return kycList
}

func (o OfferRepository) Delete(c context.Context, m model.Offer) error {
	span, c := jtrace.T().SpanFromContext(c, "OfferRepository[Delete]")
	defer span.Finish()

	if _, err := o.db.Update(c, &entity.Offer{ID: *m.ID}, persist.D{"deleted_at": time.Now()}); err != nil {
		return err
	}

	return nil
}

func (o OfferRepository) Update(c context.Context, m model.Offer) (model.Offer, error) {
	span, c := jtrace.T().SpanFromContext(c, "OfferRepository[Update]")
	defer span.Finish()

	data := mapOfferModelToEntity(m)
	offer, err := o.db.Update(c, &entity.Offer{ID: *m.ID}, data)
	if err != nil {
		return model.Offer{}, err
	}

	return mapOfferEntityToModel(*offer.(*entity.Offer)), nil
}
