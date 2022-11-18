package collection

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/contract"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
	entity "nft/internal/collection/entity"
	model "nft/internal/collection/model"

	"time"
)

type CollectionRepository struct {
	db contract.IPersist
}

type CollectionRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewCollectionRepository(params CollectionRepositoryParams) contract.ICollectionRepository {
	return &CollectionRepository{
		db: params.DB,
	}
}

func (cr CollectionRepository) Add(c context.Context, m model.Collection) (model.Collection, error) {
	span, c := jtrace.T().SpanFromContext(c, "CollectionRepository[Add]")
	defer span.Finish()

	collectionEntity := mapCollectionModelToEntity(m)
	collectionEntity.ID = uuid.New()

	createdCollection, err := cr.db.Create(c, &collectionEntity)
	if err != nil {
		return model.Collection{}, err
	}

	return mapCollectionEntityToModel(*createdCollection.(*entity.Collection)), nil
}

func (cr CollectionRepository) Update(c context.Context, m model.Collection) (model.Collection, error) {
	span, c := jtrace.T().SpanFromContext(c, "CollectionRepository[Update]")
	defer span.Finish()

	data := mapCollectionModelToEntity(m)
	updatedCollection, err := cr.db.Update(c, &entity.Collection{ID: *m.ID}, data)
	if err != nil {
		return model.Collection{}, err
	}

	return mapCollectionEntityToModel(*updatedCollection.(*entity.Collection)), nil
}

func (cr CollectionRepository) Delete(c context.Context, m model.Collection) error {
	span, c := jtrace.T().SpanFromContext(c, "CollectionRepository[Delete]")
	defer span.Finish()

	if _, err := cr.db.Update(c, &entity.Collection{ID: *m.ID, UserId: m.User.ID}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return err
	}

	return nil
}

func (cr CollectionRepository) Get(c context.Context, conditions persist.D) (model.Collection, error) {
	span, c := jtrace.T().SpanFromContext(c, "CollectionRepository[Get]")
	defer span.Finish()

	category, err := cr.db.Get(c, &entity.Collection{}, conditions)
	if err != nil {
		return model.Collection{}, err
	}

	return mapCollectionEntityToModel(*category.(*entity.Collection)), nil
}

func (cr CollectionRepository) GetAll(c context.Context, conditions persist.D) ([]model.Collection, error) {
	span, c := jtrace.T().SpanFromContext(c, "CollectionRepository[GetAll]")
	defer span.Finish()

	catList, err := cr.db.GetAll(c, &[]entity.Collection{}, conditions)
	if err != nil {
		return nil, err
	}

	return createModelCollectionListFromEntity(*catList.(*[]entity.Collection)), nil
}

func (cr CollectionRepository) HardDelete(c context.Context, id uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "CollectionRepository[HardDelete]")
	defer span.Finish()
	return cr.db.Delete(c, &entity.Collection{ID: id})
}
