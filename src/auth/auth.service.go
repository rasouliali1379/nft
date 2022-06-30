package auth

import (
	"context"
	"log"
	"maskan/client/jtrace"
	"maskan/contract"
	merrors "maskan/error"
	"maskan/pkg/crypt"
	"maskan/src/auth/model"
	jwt "maskan/src/jwt/model"

	"go.uber.org/fx"
)

type AuthService struct {
	jwtService     contract.IJwtService
	authRepository contract.IAuthRepository
	userRepository contract.IUserRepository
}

type AuthServiceParams struct {
	fx.In
	JwtService     contract.IJwtService
	AuthRepository contract.IAuthRepository
	UserRepository contract.IUserRepository
}

func NewAuthService(params AuthServiceParams) AuthService {
	return AuthService{
		jwtService:     params.JwtService,
		authRepository: params.AuthRepository,
		userRepository: params.UserRepository,
	}
}

func (a AuthService) SignUp(c context.Context, dto auth.SignUpRequest) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[SignUp]")
	defer span.Finish()

	if err := a.userRepository.EmailExists(c, dto.Email); err != nil {
		log.Println(err)
		return jwt.Jwt{}, err
	}

	if err := a.userRepository.NationalIdExists(c, dto.NationalId); err != nil {
		return jwt.Jwt{}, err
	}

	if err := a.userRepository.PhoneNumberExists(c, dto.PhoneNumber); err != nil {
		return jwt.Jwt{}, err
	}

	userId, err := a.userRepository.AddUser(c, dto)
	if err != nil {
		return jwt.Jwt{}, err
	}

	token, err := a.jwtService.Generate(c, userId)
	if err != nil {
		return jwt.Jwt{}, err
	}

	return token, nil
}

func (a AuthService) Login(c context.Context, dto auth.LoginRequest) (jwt.Jwt, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[Login]")
	defer span.Finish()

	userModel, err := a.userRepository.GetUser(c, dto.Email)
	if err != nil {
		return jwt.Jwt{}, err
	}

	newpass, _ := crypt.Hash(dto.Password)

	log.Println(newpass, userModel.Password)

	if !crypt.CompareHash(dto.Password, userModel.Password) {
		return jwt.Jwt{}, merrors.ErrInvalidCredentials
	}

	token, err := a.jwtService.Generate(c, userModel.ID.String())
	if err != nil {
		return jwt.Jwt{}, err
	}

	return token, nil
}
