package entity

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt *time.Time

	AssetId         uuid.UUID `gorm:"type:uuid;"`
	SaleId          uuid.UUID `gorm:"type:uuid;"`
	BuyerId         uuid.UUID `gorm:"type:uuid;"`
	SellerId        uuid.UUID `gorm:"type:uuid;"`
	OfferId         uuid.UUID `gorm:"type:uuid;"`
	ContractAddress string
	TransactionId   string
}
