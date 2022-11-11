package talan

import (
	dto "nft/internal/talan/dto"
	model "nft/internal/talan/model"
)

func mapAddressDtoToModel(address dto.GeneratedAddressDto) *model.Address {
	return &model.Address{
		Mnemonic:      address.Mnemonic,
		PublicAddress: address.PublicAddress,
		PrivateKey:    address.PrivateKey,
	}
}

func mapTransactionDtoToModel(transactions []dto.TransactionDto) []model.Transaction {
	var txs []model.Transaction

	for _, tx := range transactions {

		var txType model.TransactionType

		if tx.Type == "send" {
			txType = model.TransactionTypeSend
		}

		txs = append(txs, model.Transaction{
			ID:            tx.TxId,
			BlockHeight:   tx.BlockHeight,
			TimeStamp:     tx.Timestamp,
			Amount:        tx.Amount,
			Confirmations: tx.Confirmations,
			Type:          txType,
		})
	}

	return txs
}
