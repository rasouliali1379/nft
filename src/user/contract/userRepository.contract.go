package user

import (
	"context"
	model "maskan/src/auth/model"
)

type IUserRepository interface {
	NationalIdExists(c context.Context, nationalId string) error
	PhoneNumberExists(c context.Context, phoneNumber string) error
	EmailExists(c context.Context, email string) error
	AddUser(c context.Context, dto model.SignUpRequest) (model.SignUpResponse, error)
}