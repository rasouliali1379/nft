package collection

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type Collection struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	UserId      uuid.UUID `gorm:"type:uuid"`
	HeaderImage *sql.NullString
	Title       *sql.NullString
	Description *sql.NullString
	CategoryIds pq.StringArray `gorm:"type:text[]"`
	Draft       bool
}
