package auth

import (
	"context"
	"go.uber.org/fx"
	"maskan/client/jtrace"
	contract "maskan/src/auth/contract"
	model "maskan/src/auth/model"
	usercontract "maskan/src/user/contract"
)

type AuthService struct {
	authRepository contract.IAuthRepository
	userRepository usercontract.IUserRepository
}

type AuthServiceParams struct {
	fx.In
	AuthRepository contract.IAuthRepository
	UserRepository usercontract.IUserRepository
}

func NewAuthService(params AuthServiceParams) AuthService {
	return AuthService{
		authRepository: params.AuthRepository,
		userRepository: params.UserRepository,
	}
}

func (a AuthService) SignUp(c context.Context, dto model.SignUpRequest) (model.SignUpResponse, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[SignUp]")
	defer span.Finish()

	if err := a.userRepository.EmailExists(c, dto.Email); err != nil {
		return model.SignUpResponse{}, err
	}

	if err := a.userRepository.NationalIdExists(c, dto.NationalId); err != nil {
		return model.SignUpResponse{}, err
	}

	if err := a.userRepository.PhoneNumberExists(c, dto.PhoneNumber); err != nil {
		return model.SignUpResponse{}, err
	}

	response, err := a.userRepository.AddUser(c, dto)
	if err != nil {
		return model.SignUpResponse{}, err
	}

	return response, nil
}
