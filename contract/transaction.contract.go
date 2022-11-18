package contract

import (
	"context"
	"github.com/google/uuid"
	"nft/infra/persist/type"
	"nft/internal/transaction/model"
)

type ITransactionService interface {
	GetLastTransaction(c context.Context, AssetId uuid.UUID) (model.Transaction, error)
}

type ITransactionRepository interface {
	Get(c context.Context, conditions persist.D) (model.Transaction, error)
}
