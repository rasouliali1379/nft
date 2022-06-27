package auth

import "github.com/gofiber/fiber/v2"

type IAuthController interface {
	SignUp(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}
