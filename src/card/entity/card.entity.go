package card

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	UserId     uuid.UUID `gorm:"type:uuid"`
	CardNumber string
	IBAN       string
	ApprovedBy *uuid.UUID `gorm:"type:uuid"`
}
