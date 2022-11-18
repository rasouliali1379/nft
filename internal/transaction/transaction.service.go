package transaction

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/contract"
	"nft/infra/jtrace"
	"nft/infra/persist/type"
	model "nft/internal/transaction/model"
)

type TransactionService struct {
	saleRepository contract.ITransactionRepository
}

type TransactionServiceParams struct {
	fx.In
	TransactionRepository contract.ITransactionRepository
}

func NewTransactionService(params TransactionServiceParams) contract.ITransactionService {
	return TransactionService{
		saleRepository: params.TransactionRepository,
	}
}

func (t TransactionService) GetLastTransaction(c context.Context, AssetId uuid.UUID) (model.Transaction, error) {
	span, c := jtrace.T().SpanFromContext(c, "TransactionService[GetLastTransaction]")
	defer span.Finish()
	return t.saleRepository.Get(c, persist.D{"asset_id": AssetId})
}
