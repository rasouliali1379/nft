package auth

import (
	"context"
	"log"
	"maskan/client/jtrace"
	"maskan/contract"
	merrors "maskan/error"
	"maskan/pkg/crypt"
	jwt "maskan/src/jwt/model"
	user "maskan/src/user/model"

	"go.uber.org/fx"
)

type AuthService struct {
	jwtService     contract.IJwtService
	userService    contract.IUserService
	authRepository contract.IAuthRepository
}

type AuthServiceParams struct {
	fx.In
	JwtService     contract.IJwtService
	UserService    contract.IUserService
	AuthRepository contract.IAuthRepository
}

func NewAuthService(params AuthServiceParams) AuthService {
	return AuthService{
		jwtService:     params.JwtService,
		authRepository: params.AuthRepository,
		userService:    params.UserService,
	}
}

func (a AuthService) SignUp(c context.Context, userModel user.User) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[SignUp]")
	defer span.Finish()

	createdUser, err := a.userService.AddUser(c, userModel)
	if err != nil {
		return jwt.Jwt{}, err
	}

	token, err := a.jwtService.Generate(c, createdUser.ID.String())
	if err != nil {
		return jwt.Jwt{}, err
	}

	return token, nil
}

func (a AuthService) Login(c context.Context, email string, password string) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[Login]")
	defer span.Finish()

	userModel, err := a.userService.GetUser(c, user.UserQuery{Email: email})
	if err != nil {
		return jwt.Jwt{}, err
	}

	log.Println(userModel)

	if !crypt.CompareHash(password, userModel.Password) {
		return jwt.Jwt{}, merrors.ErrInvalidCredentials
	}

	token, err := a.jwtService.Generate(c, userModel.ID.String())
	if err != nil {
		return jwt.Jwt{}, err
	}

	return token, nil
}
