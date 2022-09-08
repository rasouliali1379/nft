package nft

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/client/jtrace"
	"nft/client/persist/model"
	"nft/contract"
	entity "nft/src/nft/entity"
	model "nft/src/nft/model"
	"time"
)

type NftRepository struct {
	db contract.IPersist
}

type NftRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewNftRepository(params NftRepositoryParams) contract.INftRepository {
	return &NftRepository{
		db: params.DB,
	}
}

func (n NftRepository) Exists(c context.Context, conditions persist.Conds) error {
	//TODO implement me
	panic("implement me")
}

func (n NftRepository) Add(c context.Context, m model.Nft) (model.Nft, error) {
	span, c := jtrace.T().SpanFromContext(c, "NftRepository[Add]")
	defer span.Finish()

	nftEntity := mapNftModelToEntity(m)
	nftEntity.ID = uuid.New()

	createdNft, err := n.db.Create(c, &nftEntity)
	if err != nil {
		return model.Nft{}, err
	}

	return mapNftEntityToModel(*createdNft.(*entity.Nft)), nil
}

func (n NftRepository) Update(c context.Context, m model.Nft) (model.Nft, error) {
	span, c := jtrace.T().SpanFromContext(c, "NftRepository[Update]")
	defer span.Finish()

	data := mapNftModelToEntity(m)
	updatedNft, err := n.db.Update(c, &entity.Nft{ID: *m.ID}, data)
	if err != nil {
		return model.Nft{}, err
	}

	return mapNftEntityToModel(*updatedNft.(*entity.Nft)), nil
}

func (n NftRepository) Delete(c context.Context, id uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "NftRepository[Delete]")
	defer span.Finish()

	if _, err := n.db.Update(c, &entity.Nft{ID: id}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return err
	}

	return nil
}
func (n NftRepository) HardDelete(c context.Context, id uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "NftRepository[HardDelete]")
	defer span.Finish()
	return n.db.Delete(c, &entity.Nft{ID: id})
}

func (n NftRepository) Get(c context.Context, conditions persist.Conds) (model.Nft, error) {
	span, c := jtrace.T().SpanFromContext(c, "NftRepository[Get]")
	defer span.Finish()

	category, err := n.db.Get(c, &entity.Nft{}, conditions)
	if err != nil {
		return model.Nft{}, err
	}

	return mapNftEntityToModel(*category.(*entity.Nft)), nil
}

func (n NftRepository) GetAll(c context.Context, conditions persist.Conds) ([]model.Nft, error) {
	span, c := jtrace.T().SpanFromContext(c, "NftRepository[GetAll]")
	defer span.Finish()

	catList, err := n.db.GetAll(c, &[]entity.Nft{}, conditions)
	if err != nil {
		return nil, err
	}

	return createModelNftListFromEntity(*catList.(*[]entity.Nft)), nil
}
