package user

import (
	"context"
	"maskan/client/jtrace"
	"maskan/contract"

	model "maskan/src/user/model"

	"go.uber.org/fx"
)

type UserService struct {
	jwtService     contract.IJwtService
	userRepository contract.IUserRepository
}

type UserServiceParams struct {
	fx.In
	JwtService     contract.IJwtService
	UserRepository contract.IUserRepository
}

func NewUserService(params UserServiceParams) UserService {
	return UserService{
		jwtService:     params.JwtService,
		userRepository: params.UserRepository,
	}
}

func (u UserService) GetAllUsers(c context.Context) ([]model.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[GetAllUsers]")
	defer span.Finish()

	userList, err := u.userRepository.GetAllUsers(c)
	if err != nil {
		return nil, err
	}

	return userList, nil
}

func (u UserService) GetUser(c context.Context, query model.UserQuery) (model.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[GetUser]")
	defer span.Finish()
	return u.userRepository.GetUser(c, query)
}

func (u UserService) AddUser(c context.Context, userModel model.User) (model.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[AddUser]")
	defer span.Finish()

	if err := u.userRepository.EmailExists(c, userModel.Email); err != nil {
		return model.User{}, err
	}

	if err := u.userRepository.NationalIdExists(c, userModel.NationalId); err != nil {
		return model.User{}, err
	}

	if err := u.userRepository.PhoneNumberExists(c, userModel.PhoneNumber); err != nil {
		return model.User{}, err
	}

	newUser, err := u.userRepository.AddUser(c, userModel)
	if err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (u UserService) UpdateUser(c context.Context, userModel model.User) (model.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "service[UpdateUser]")
	defer span.Finish()

	userRecord, err := u.GetUser(c, model.UserQuery{
		ID: userModel.ID.String(),
	})

	if err != nil {
		return model.User{}, err
	}

	if userRecord.Email != userModel.Email && len(userModel.Email) > 0 {
		if err := u.userRepository.EmailExists(c, userModel.Email); err != nil {
			return model.User{}, err
		}
	}

	if userRecord.NationalId != userModel.NationalId && len(userModel.NationalId) > 0 {
		if err := u.userRepository.NationalIdExists(c, userModel.NationalId); err != nil {
			return model.User{}, err
		}
	}

	if userRecord.PhoneNumber != userModel.PhoneNumber && len(userModel.PhoneNumber) > 0 {
		if err := u.userRepository.PhoneNumberExists(c, userModel.PhoneNumber); err != nil {
			return model.User{}, err
		}
	}

	_, err = u.userRepository.UpdateUser(c, userModel)
	if err != nil {
		return model.User{}, err
	}

	userRecord, err = u.GetUser(c, model.UserQuery{
		ID: userModel.ID.String(),
	})

	if err != nil {
		return model.User{}, err
	}

	return userRecord, nil
}

func (u UserService) DeleteUser(c context.Context, userId string) error {
	span, c := jtrace.T().SpanFromContext(c, "service[DeleteUser]")
	defer span.Finish()

	if _, err := u.GetUser(c, model.UserQuery{ID: userId}); err != nil {
		return err
	}

	return u.userRepository.DeleteUser(c, userId)
}
