package kyc

import (
	"context"
	"nft/contract"
	"nft/infra/jtrace"
	"time"

	persist "nft/infra/persist/model"
	entity "nft/internal/kyc/entity"
	model "nft/internal/kyc/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type KycRepository struct {
	db contract.IPersist
}

type KycRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewKYCRepository(params KycRepositoryParams) contract.IKycRepository {
	return &KycRepository{
		db: params.DB,
	}
}

func (k KycRepository) Exists(c context.Context, conditions persist.Conds) error {
	span, c := jtrace.T().SpanFromContext(c, "KycRepository[Exists]")
	defer span.Finish()

	if _, err := k.db.Get(c, &entity.Kyc{}, conditions); err != nil {
		return err
	}

	return nil
}

func (k KycRepository) Add(c context.Context, kyc model.Kyc) (model.Kyc, error) {
	span, c := jtrace.T().SpanFromContext(c, "KycRepository[Add]")
	defer span.Finish()

	kycEntity := mapKycModelToEntity(kyc)
	kycEntity.ID = uuid.New()

	appeal, err := k.db.Create(c, &kycEntity)
	if err != nil {
		return model.Kyc{}, err
	}

	return mapKycEntityToModel(appeal.(*entity.Kyc)), nil
}

func (k KycRepository) Update(c context.Context, kyc model.Kyc) (model.Kyc, error) {
	span, c := jtrace.T().SpanFromContext(c, "KycRepository[Update]")
	defer span.Finish()

	data := mapKycModelToEntity(kyc)
	updatedAppeal, err := k.db.Update(c, &entity.Kyc{ID: kyc.ID}, data)
	if err != nil {
		return model.Kyc{}, err
	}

	return mapKycEntityToModel(updatedAppeal.(*entity.Kyc)), nil
}

func (k KycRepository) Delete(c context.Context, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "KycRepository[Delete]")
	defer span.Finish()

	if _, err := k.db.Update(c, &entity.Kyc{ID: userId}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return err
	}

	return nil
}

func (k KycRepository) Get(c context.Context, conditions persist.Conds) (model.Kyc, error) {
	span, c := jtrace.T().SpanFromContext(c, "KycRepository[Get]")
	defer span.Finish()

	category, err := k.db.Get(c, &entity.Kyc{}, conditions)
	if err != nil {
		return model.Kyc{}, err
	}

	return mapKycEntityToModel(category.(*entity.Kyc)), nil
}

func (k KycRepository) GetAll(c context.Context, conditions persist.Conds) ([]model.Kyc, error) {
	span, c := jtrace.T().SpanFromContext(c, "KycRepository[GetAll]")
	defer span.Finish()

	catList, err := k.db.GetAll(c, &[]entity.Kyc{}, conditions)
	if err != nil {
		return nil, err
	}

	return createModelKycList(catList.(*[]entity.Kyc)), nil
}
