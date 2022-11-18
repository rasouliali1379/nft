package transaction

import (
	"nft/internal/transaction/entity"
	"nft/internal/transaction/model"
)

func mapTransactionEntityToModel(e entity.Transaction) model.Transaction {
	return model.Transaction{
		AssetId:         e.AssetId,
		SaleId:          e.SaleId,
		BuyerId:         e.BuyerId,
		SellerId:        e.SellerId,
		OfferId:         e.OfferId,
		ContractAddress: e.ContractAddress,
		TransactionId:   e.TransactionId,
	}
}
