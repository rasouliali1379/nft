package category

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	Name      string
	ParentId  *uuid.UUID `gorm:"type:uuid;"`
	CreatedBy uuid.UUID  `gorm:"type:uuid;"`
}
