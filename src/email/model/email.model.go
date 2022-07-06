package email

import (
	"time"

	"github.com/google/uuid"
)

type Email struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt *time.Time

	UserId   uuid.UUID
	OtpCode  string
	Verified bool
}
