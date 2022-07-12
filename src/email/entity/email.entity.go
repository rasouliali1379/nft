package email

import (
	"time"

	"github.com/google/uuid"
)

type Email struct {
	ID        uint `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	UserId   uuid.UUID `gorm:"type:uuid;"`
	Email    string    `gorm:"not null"`
	Verified bool      `gorm:"default:false"`
}
