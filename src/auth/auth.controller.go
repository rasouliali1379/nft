package auth

import (
	"errors"
	"maskan/client/jtrace"
	"maskan/pkg/filper"
	contract "maskan/src/auth/contract"
	model "maskan/src/auth/model"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService contract.IAuthService
}

func NewAuthController(service contract.IAuthService) contract.IAuthController {
	return AuthController{
		authService: service,
	}
}

func (a AuthController) SignUp(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller")
	defer span.Finish()

	var dto model.SignUpRequest
	if err := c.BodyParser(&dto); err != nil {
		return err
	}

	//validation

	response, err := a.authService.SignUp(ctx, dto)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			return filper.GetBadRequestError(c, "email already exists")
		}

		if errors.Is(err, ErrPhoneNumberExists) {
			return filper.GetBadRequestError(c, "phone number already exists")
		}

		if errors.Is(err, ErrNationalIdExists) {
			return filper.GetBadRequestError(c, "national id already exists")
		}

		return fiber.ErrInternalServerError
	}

	return c.JSON(response)
}

func (a AuthController) Login(c *fiber.Ctx) error {
	span, _ := jtrace.T().SpanFromContext(c.Context(), "controller")
	defer span.Finish()
	return c.SendString("Hello, World!")
}
