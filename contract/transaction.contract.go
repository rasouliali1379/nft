package contract

import (
	"context"
	"github.com/google/uuid"
	persist "nft/infra/persist/model"
	"nft/internal/transaction/model"
)

type ITransactionService interface {
	GetLastTransaction(c context.Context, AssetId uuid.UUID) (model.Transaction, error)
}

type ITransactionRepository interface {
	Get(c context.Context, conditions persist.Conds) (model.Transaction, error)
}
