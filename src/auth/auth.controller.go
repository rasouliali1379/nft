package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"log"
	"maskan/client/jtrace"
	"maskan/contract"
	merror "maskan/error"
	"maskan/pkg/filper"
	"maskan/pkg/validator"
	dto "maskan/src/auth/dto"
	jwt "maskan/src/jwt/model"
	user "maskan/src/user"
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

// SignUp godoc
// @Summary  sign up new user
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    message  body      dto.SignUpRequest  true  "sign up request body"
// @Success  201      {object}  jwt.Jwt
// @Router   /v1/auth/signup [post]
func (a AuthController) SignUp(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[SignUp]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var dto dto.SignUpRequest
	if err := c.BodyParser(&dto); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(dto)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)

	}
	
	var response jwt.Jwt
	response, err := a.authService.SignUp(ctx, user.MapSignUpDtoToUserModel(dto, ""))
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

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (a AuthController) Login(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[Login]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var dto dto.LoginRequest
	if err := c.BodyParser(&dto); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(dto)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)

	}

	response, err := a.authService.Login(ctx, dto.Email, dto.Password)
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

	var dto dto.RefreshRequest
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
