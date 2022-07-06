package contract

import (
	"context"
	model "maskan/src/user/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	DeleteUser(c context.Context, userId uuid.UUID) error
}

type IUserRepository interface {
	UserExists(c context.Context, query model.UserQuery) error
	AddUser(c context.Context, user model.User) (model.User, error)
	UpdateUser(c context.Context, userModel model.User) (model.User, error)
	DeleteUser(c context.Context, userId uuid.UUID) error
	GetUser(c context.Context, query model.UserQuery) (model.User, error)
	GetAllUsers(c context.Context) ([]model.User, error)
}
