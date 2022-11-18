package transaction

import (
	"context"
	"go.uber.org/fx"
	"nft/contract"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
	entity "nft/internal/transaction/entity"
	"nft/internal/transaction/model"
)

type TransactionRepository struct {
	db contract.IPersist
}

type TransactionRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewTransactionRepository(params TransactionRepositoryParams) contract.ITransactionRepository {
	return &TransactionRepository{
		db: params.DB,
	}
}

func (t TransactionRepository) Get(c context.Context, conditions persist.D) (model.Transaction, error) {
	span, c := jtrace.T().SpanFromContext(c, "TransactionRepository[Get]")
	defer span.Finish()

	tx, err := t.db.Get(c, &entity.Transaction{}, conditions)
	if err != nil {
		return model.Transaction{}, err
	}

	return mapTransactionEntityToModel(*tx.(*entity.Transaction)), nil
}
