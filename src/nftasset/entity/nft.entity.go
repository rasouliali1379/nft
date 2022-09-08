package nftasset

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type Nft struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	ApprovedBy      *uuid.UUID `gorm:"type:uuid"`
	RejectedBy      *uuid.UUID `gorm:"type:uuid"`
	UserId          uuid.UUID  `gorm:"type:uuid"`
	RejectionReason *sql.NullString
	NftImage        *sql.NullString
	Title           *sql.NullString
	Description     *sql.NullString
	CategoryIds     pq.StringArray `gorm:"type:text[]"`
	Draft           bool
}
