package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	NationalId     string `gorm:"unique"`
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
