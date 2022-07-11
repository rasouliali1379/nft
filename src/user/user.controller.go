package user

import (
	"maskan/client/jtrace"
	"maskan/contract"
	"maskan/pkg/filper"
	"maskan/pkg/validator"
	authdto "maskan/src/auth/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type UserController struct {
	authService contract.IAuthService
	jwtService  contract.IJwtService
	userService contract.IUserService
}

type UserControllerParams struct {
	fx.In
	AuthService contract.IAuthService
	JwtService  contract.IJwtService
	UserService contract.IUserService
}

func NewUserController(params UserControllerParams) contract.IUserController {
	return UserController{
		authService: params.AuthService,
		jwtService:  params.JwtService,
		userService: params.UserService,
	}
}

func (u UserController) GetAllUsers(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[GetAllUsers]")
	defer span.Finish()

	users, err := u.userService.GetAllUsers(ctx)
	if err != nil {
		return filper.GetInternalError(c, "")
	}
	return c.JSON(createUserList(users))
}

func (u UserController) GetUser(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[GetUser]")
	defer span.Finish()

	userId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid user id")
	}

	userModel, err := u.userService.GetUser(ctx, map[string]any{"id": userId})
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.JSON(mapUserModelToResponse(userModel))
}

func (u UserController) AddUser(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[AddUser]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var dto authdto.SignUpRequest
	if err := c.BodyParser(&dto); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errs := validator.Validate(dto)
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)

	}

	userModel, err := u.userService.AddUser(ctx, MapSignUpDtoToUserModel(dto, uuid.UUID{}))
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapUserModelToResponse(userModel))
}

func (u UserController) UpdateUser(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[UpdateUser]")
	defer span.Finish()

	userId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid user id")
	}

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var dto authdto.SignUpRequest
	if err := c.BodyParser(&dto); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	userModel, err := u.userService.UpdateUser(ctx, MapSignUpDtoToUserModel(dto, userId))
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.JSON(mapUserModelToResponse(userModel))
}
func (u UserController) DeleteUser(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "controller[DeleteUser]")
	defer span.Finish()

	userId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid user id")
	}

	err = u.userService.DeleteUser(ctx, userId)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "user deleted successfully")
}
