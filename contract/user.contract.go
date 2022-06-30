package contract

import (
	"context"
	model "maskan/src/auth/model"
	user "maskan/src/user/model"
)

type IUserRepository interface {
	NationalIdExists(c context.Context, nationalId string) error
	PhoneNumberExists(c context.Context, phoneNumber string) error
	EmailExists(c context.Context, email string) error
	AddUser(c context.Context, dto model.SignUpRequest) (string, error)
	GetUser(c context.Context, email string) (user.User, error)
}
