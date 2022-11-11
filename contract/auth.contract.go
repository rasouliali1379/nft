package contract

import (
	"context"

	jwt "nft/internal/jwt/model"
	user "nft/internal/user/model"

	"github.com/gofiber/fiber/v2"
)

type IAuthController interface {
	SignUp(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	ResendEmail(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type IAuthService interface {
	SignUp(c context.Context, model user.User) (string, error)
	Login(c context.Context, email string, password string) (jwt.Jwt, error)
	VerifyEmail(c context.Context, token string, code string) (jwt.Jwt, error)
	ResendVerificationEmail(c context.Context, token string) (string, error)
}
