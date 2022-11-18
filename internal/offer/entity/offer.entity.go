package entity

import (
	"github.com/google/uuid"
	"time"
)

type Offer struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	DeletedAt *time.Time

	UserId   uuid.UUID `gorm:"type:uuid;"`
	SaleId   uuid.UUID `gorm:"type:uuid;"`
	Price    float64
	Accepted bool
}
