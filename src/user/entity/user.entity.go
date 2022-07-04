package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time 

	NationalId     string
	FirstName      string
	LastName       string
	Email          string 
	PhoneNumber    string 
	Password       string
	LandLineNumber string
	Province       string
	City           string
	Address        string
	PublicKey      string
	PrivateKey     string
}
