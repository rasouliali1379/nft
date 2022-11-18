package contract

import (
	"context"
	"nft/infra/persist/type"
	model "nft/internal/user/model"

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
	GetUser(c context.Context, conditions persist.D) (model.User, error)
	AddUser(c context.Context, userModel model.User) (model.User, error)
	UpdateUser(c context.Context, userModel model.User) (model.User, error)
	DeleteUser(c context.Context, userId uuid.UUID) error
}

type IUserRepository interface {
	Exists(c context.Context, conditions persist.D) (bool, error)
	Add(c context.Context, user model.User) (model.User, error)
	Update(c context.Context, userModel model.User) (model.User, error)
	Delete(c context.Context, userId uuid.UUID) error
	Get(c context.Context, conditions persist.D) (model.User, error)
	GetAll(c context.Context) ([]model.User, error)
}
