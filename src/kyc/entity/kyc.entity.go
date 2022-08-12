package entity

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type KYC struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	ApprovedBy      *uuid.UUID `gorm:"type:uuid"`
	RejectedBy      *uuid.UUID `gorm:"type:uuid"`
	UserId          uuid.UUID  `gorm:"type:uuid"`
	RejectionReason *sql.NullString
	IdCardImage     string
	PortraitImage   string
}
