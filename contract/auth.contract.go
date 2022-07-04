package contract

import (
	"context"

	"github.com/gofiber/fiber/v2"
	jwt "maskan/src/jwt/model"
	user "maskan/src/user/model"
)

type IAuthController interface {
	SignUp(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type IAuthService interface {
	SignUp(c context.Context, model user.User) (jwt.Jwt, error)
	Login(c context.Context, email string, password string) (jwt.Jwt, error)
}

type IAuthRepository interface{}
