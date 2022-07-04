package contract

import (
	"context"
	"github.com/gofiber/fiber/v2"
	model "maskan/src/user/model"
)

type IUserController interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	AddUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type IUserService interface {
	GetAllUsers(c context.Context) ([]model.User, error)
	GetUser(c context.Context, query model.UserQuery) (model.User, error)
	AddUser(c context.Context, userModel model.User) (model.User, error)
	UpdateUser(c context.Context, userModel model.User) (model.User, error)
	DeleteUser(c context.Context, userId string) error
}

type IUserRepository interface {
	NationalIdExists(c context.Context, nationalId string) error
	PhoneNumberExists(c context.Context, phoneNumber string) error
	EmailExists(c context.Context, email string) error
	AddUser(c context.Context, user model.User) (model.User, error)
	UpdateUser(c context.Context, userModel model.User) (model.User, error)
	DeleteUser(c context.Context, userId string) error
	GetUser(c context.Context, query model.UserQuery) (model.User, error)
	GetAllUsers(c context.Context) ([]model.User, error)
}
