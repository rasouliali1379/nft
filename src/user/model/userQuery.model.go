package user

import "github.com/google/uuid"

type UserQuery struct {
	ID          uuid.UUID
	Email       string
	NationalId  string
	PhoneNumber string
}
