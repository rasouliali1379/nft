package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	FirstName      string
	LastName       string
	NationalId     string
	Email          string
	PhoneNumber    string
	Password       string
	LandLineNumber string
	Province       string
	City           string
	Address        string
	PublicKey      string
	PrivateKey     string
	Mnemonic       string
}
