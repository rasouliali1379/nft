package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

var (
	Bcrypt Micro = &micro{}
)

type Micro interface {
	Hash(password string) (string, error)
	CompareHash(password, hash string) bool
}

type micro struct{}

func (m *micro) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (m *micro) CompareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
