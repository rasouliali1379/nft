package auth

import (
	"errors"
	"log"
	"maskan/client/jtrace"
	"maskan/contract"
	merror "maskan/error"
	"maskan/pkg/filper"
	"maskan/pkg/validator"
	model "maskan/src/auth/model"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type AuthController struct {
	authService contract.IAuthService
	jwtService  contract.IJwtService
}

type AuthControllerParams struct {
	fx.In
	AuthService contract.IAuthService
	JwtService  contract.IJwtService
}

func NewAuthController(params AuthControllerParams) contract.IAuthController {
	return AuthController{
		authService: params.AuthService,
		jwtService:  params.JwtService,
	}
}

func (a AuthController) SignUp(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[SignUp]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var dto model.SignUpRequest
	if err := c.BodyParser(&dto); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(dto)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)

	}

	response, err := a.authService.SignUp(ctx, dto)
	if err != nil {

		if errors.Is(err, merror.ErrEmailExists) {
			return filper.GetBadRequestError(c, "email already exists")
		}

		if errors.Is(err, merror.ErrPhoneNumberExists) {
			return filper.GetBadRequestError(c, "phone number already exists")
		}

		if errors.Is(err, merror.ErrNationalIdExists) {
			return filper.GetBadRequestError(c, "national id already exists")
		}

		return filper.GetInternalError(c, "")
	}

	return c.JSON(response)
}

func (a AuthController) Login(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[Login]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var dto model.LoginRequest
	if err := c.BodyParser(&dto); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(dto)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)

	}

	response, err := a.authService.Login(ctx, dto)
	if err != nil {
		if errors.Is(err, merror.ErrInvalidCredentials) {
			return filper.GetInvalidCredentialsError(c, "invalid credentials")
		}
		log.Println(err)
		return filper.GetInternalError(c, "")
	}

	return c.JSON(response)
}

func (a AuthController) Refresh(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[Refresh]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var dto model.RefreshRequest
	if err := c.BodyParser(&dto); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(dto)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)

	}

	response, err := a.jwtService.Refresh(ctx, dto.RefreshToken)
	if err != nil {
		if errors.Is(err, merror.ErrInvalidCredentials) {
			return filper.GetInvalidCredentialsError(c, "invalid credentials")
		}
		log.Println(err)
		return filper.GetInternalError(c, "")
	}

	return c.JSON(response)
}
