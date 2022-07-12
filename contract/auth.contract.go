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
	VerifyEmail(c *fiber.Ctx) error	
	ResendEmail(c *fiber.Ctx) error
}

type IAuthService interface {
	SignUp(c context.Context, model user.User) (string, error)
	Login(c context.Context, email string, password string) (jwt.Jwt, error)
	VerifyEmail(c context.Context, token string, code string) (jwt.Jwt, error)
	ResendVerificationEmail(c context.Context, token string) (string, error)
}