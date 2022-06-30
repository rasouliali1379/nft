package contract

import (
	"context"

	"github.com/gofiber/fiber/v2"
	auth "maskan/src/auth/model"
	jwt "maskan/src/jwt/model"
)

type IAuthController interface {
	SignUp(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type IAuthService interface {
	SignUp(context.Context, auth.SignUpRequest) (jwt.Jwt, error)
	Login(context.Context, auth.LoginRequest) (jwt.Jwt, error)
}

type IAuthRepository interface{}
