package user

import (
	"nft/contract"
	"nft/infra/jtrace"
	authdto "nft/internal/auth/dto"
	user "nft/internal/user/dto"
	"nft/pkg/filper"
	"nft/pkg/validator"

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

// GetAllUsers godoc
// @Summary  get users list
// @Tags     user
// @Accept   json
// @Produce  json
// @Success  200      {object}  user.UserList
// @Router   /v1/user [get]
func (u UserController) GetAllUsers(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "UserController[GetAllUsers]")
	defer span.Finish()

	users, err := u.userService.GetAllUsers(ctx)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	var userList user.UserList = createUserList(users)
	return c.JSON(userList)
}

// GetUser godoc
// @Summary  get single user
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    id   path      int  true  "user id"
// @Success  200  {object}  user.User
// @Router   /v1/user/{id} [get]
func (u UserController) GetUser(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "UserController[GetUser]")
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

// AddUser godoc
// @Summary  add user
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    message  body      authdto.SignUpRequest  true  "add user request body"
// @Success  200      {object}  user.User
// @Router   /v1/user [post]
func (u UserController) AddUser(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "UserController[AddUser]")
	defer span.Finish()

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var dto authdto.SignUpRequest
	if err := c.BodyParser(&dto); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errRes := validator.Validate(dto)
	if len(errRes.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)

	}

	userModel, err := u.userService.AddUser(ctx, MapSignUpDtoToUserModel(dto, uuid.UUID{}))
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapUserModelToResponse(userModel))
}

// UpdateUser godoc
// @Summary  update existing user
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    id       path      int                    true  "user id that will be updated"
// @Param    message  body      authdto.SignUpRequest  true  "update user request body"
// @Success  200  {object}  user.User
// @Router   /v1/user/{id} [patch]
func (u UserController) UpdateUser(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "UserController[UpdateUser]")
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

// DeleteUser godoc
// @Summary  delete existing user
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "user id that will be deleted"
// @Success  200  {string}  string  "user deleted successfully"
// @Router   /v1/user/{id} [delete]
func (u UserController) DeleteUser(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "UserController[DeleteUser]")
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
