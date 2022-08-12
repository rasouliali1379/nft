package kyc

import (
	"context"
	"nft/client/jtrace"
	"nft/contract"
	"time"

	persist "nft/client/persist/model"
	entity "nft/src/kyc/entity"
	model "nft/src/kyc/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type KYCRepository struct {
	db contract.IPersist
}

type KYCRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewKYCRepository(params KYCRepositoryParams) contract.IKYCRepository {
	return &KYCRepository{
		db: params.DB,
	}
}

func (k KYCRepository) Exists(c context.Context, conditions persist.Conds) error {
	span, c := jtrace.T().SpanFromContext(c, "KYCRepository[Exists]")
	defer span.Finish()

	if _, err := k.db.Get(c, &entity.KYC{}, conditions); err != nil {
		return err
	}

	return nil
}

func (k KYCRepository) Add(c context.Context, kyc model.KYC) (model.KYC, error) {
	span, c := jtrace.T().SpanFromContext(c, "KYCRepository[Add]")
	defer span.Finish()

	kycEntity := mapKYCModelToEntity(kyc)
	kycEntity.ID = uuid.New()

	appeal, err := k.db.Create(c, &kycEntity)
	if err != nil {
		return model.KYC{}, err
	}

	return mapKYCEntityToModel(appeal.(*entity.KYC)), nil
}

func (k KYCRepository) Update(c context.Context, kyc model.KYC) (model.KYC, error) {
	span, c := jtrace.T().SpanFromContext(c, "KYCRepository[Update]")
	defer span.Finish()

	data := mapKYCModelToEntity(kyc)
	updatedAppeal, err := k.db.Update(c, &entity.KYC{ID: kyc.ID}, data)
	if err != nil {
		return model.KYC{}, err
	}

	return mapKYCEntityToModel(updatedAppeal.(*entity.KYC)), nil
}

func (k KYCRepository) Delete(c context.Context, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "KYCRepository[Delete]")
	defer span.Finish()

	if _, err := k.db.Update(c, &entity.KYC{ID: userId}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return err
	}

	return nil
}

func (k KYCRepository) Get(c context.Context, conditions persist.Conds) (model.KYC, error) {
	span, c := jtrace.T().SpanFromContext(c, "KYCRepository[Get]")
	defer span.Finish()

	category, err := k.db.Get(c, &entity.KYC{}, conditions)
	if err != nil {
		return model.KYC{}, err
	}

	return mapKYCEntityToModel(category.(*entity.KYC)), nil
}

func (k KYCRepository) GetAll(c context.Context, conditions persist.Conds) ([]model.KYC, error) {
	span, c := jtrace.T().SpanFromContext(c, "KYCRepository[GetAll]")
	defer span.Finish()

	catList, err := k.db.GetAll(c, &[]entity.KYC{}, conditions)
	if err != nil {
		return nil, err
	}

	return createModelKYCList(catList.(*[]entity.KYC)), nil
}
