package auth

import (
	"errors"
	"log"

	"nft/client/jtrace"
	"nft/contract"
	merror "nft/error"
	"nft/pkg/filper"
	"nft/pkg/validator"
	dto "nft/src/auth/dto"
	jwt "nft/src/jwt/model"
	user "nft/src/user"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	return &AuthController{
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
// @Success  201      {object}  dto.OtpToken
// @Router   /v1/auth/signup [post]
func (a AuthController) SignUp(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "AuthController[SignUp]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var signUpRequest dto.SignUpRequest
	if err := c.BodyParser(&signUpRequest); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(signUpRequest)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)

	}

	token, err := a.authService.SignUp(ctx, user.MapSignUpDtoToUserModel(signUpRequest, uuid.UUID{}))
	if err != nil {
		log.Println(err)
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

	return c.Status(fiber.StatusCreated).JSON(dto.OtpToken{Token: token})
}

// Login godoc
// @Summary  login user
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    message  body      dto.LoginRequest  true  "login request body"
// @Success  200      {object}  jwt.Jwt
// @Router   /v1/auth/login [post]
func (a AuthController) Login(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "AuthController[Login]")
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

	var response jwt.Jwt
	response, err := a.authService.Login(ctx, dto.Email, dto.Password)
	if err != nil {
		if errors.Is(err, merror.ErrInvalidCredentials) {
			return filper.GetUnAuthError(c, "invalid credentials")
		}
		return filper.GetInternalError(c, "")
	}

	return c.JSON(response)
}

// Refresh godoc
// @Summary  refresh user token
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    message  body      dto.RefreshRequest  true  "refresh token request body"
// @Success  200      {object}  jwt.Jwt
// @Router   /v1/auth/refresh [post]
func (a AuthController) Refresh(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "AuthController[Refresh]")
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
		if errors.Is(err, merror.ErrTokenInvoked) {
			return filper.GetUnAuthError(c, "token invoked")
		}
		return filper.GetInternalError(c, "")
	}

	return c.JSON(response)
}

// VerifyEmail godoc
// @Summary  verify user email
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    message  body      dto.VerifyEmailRequest  true  "verify email request body"
// @Success  200      {object}  jwt.Jwt
// @Router   /v1/auth/verify-email [post]
func (a AuthController) VerifyEmail(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "AuthController[VerifyEmail]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.VerifyEmailRequest
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(request)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	response, err := a.authService.VerifyEmail(ctx, request.Token, request.Code)
	if err != nil {
		if errors.Is(err, merror.ErrInvalidCredentials) {
			return filper.GetUnAuthError(c, "invalid credentials")
		}
		return filper.GetInternalError(c, "")
	}

	return c.JSON(response)
}

// ResendEmail godoc
// @Summary  resend verification email
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    message  body      dto.ResendEmailRequest  true  "resend email request body"
// @Success  200      {object}  dto.OtpToken
// @Router   /v1/auth/resend-email [post]
func (a AuthController) ResendEmail(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "AuthController[ResendEmail]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.ResendEmailRequest
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(request)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	token, err := a.authService.ResendVerificationEmail(ctx, request.Token)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.JSON(dto.OtpToken{Token: token})
}
