package model

import "github.com/google/uuid"

type Transaction struct {
	AssetId         uuid.UUID `gorm:"type:uuid;"`
	SaleId          uuid.UUID `gorm:"type:uuid;"`
	BuyerId         uuid.UUID `gorm:"type:uuid;"`
	SellerId        uuid.UUID `gorm:"type:uuid;"`
	OfferId         uuid.UUID `gorm:"type:uuid;"`
	ContractAddress string
	TransactionId   string
}
