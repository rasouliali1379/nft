package user

import (
	"context"
	"errors"
	"nft/client/jtrace"
	"nft/contract"

	merror "nft/error"
	model "nft/src/user/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type UserService struct {
	jwtService     contract.IJwtService
	userRepository contract.IUserRepository
	emailService   contract.IEmailService
}

type UserServiceParams struct {
	fx.In
	JwtService     contract.IJwtService
	UserRepository contract.IUserRepository
	EmailService   contract.IEmailService
}

func NewUserService(params UserServiceParams) UserService {
	return UserService{
		jwtService:     params.JwtService,
		userRepository: params.UserRepository,
		emailService:   params.EmailService,
	}
}

func (u UserService) GetAllUsers(c context.Context) ([]model.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserService[GetAllUsers]")
	defer span.Finish()

	userList, err := u.userRepository.GetAll(c)
	if err != nil {
		return nil, err
	}

	for i, item := range userList {
		userEmail, err := u.emailService.GetLastVerifiedEmail(c, item.ID)
		if err != nil {
			if errors.Is(err, merror.ErrRecordNotFound) {
				continue
			}
			return nil, err
		}

		userList[i].Email = userEmail.Email
	}

	return userList, nil
}

func (u UserService) GetUser(c context.Context, conditions map[string]any) (model.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserService[GetUser]")
	defer span.Finish()

	userModel, err := u.userRepository.Get(c, conditions)
	if err != nil {
		return model.User{}, err
	}

	userEmail, err := u.emailService.GetLastVerifiedEmail(c, userModel.ID)
	if err != nil {
		return model.User{}, err
	}

	userModel.Email = userEmail.Email
	return userModel, nil
}

func (u UserService) AddUser(c context.Context, userModel model.User) (model.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserService[AddUser]")
	defer span.Finish()

	if err := u.emailService.EmailExists(c, userModel.Email); err != nil {
		return model.User{}, err
	}

	if err := u.userRepository.Exists(c, map[string]any{"national_id": userModel.NationalId}); err != nil {
		return model.User{}, err
	}

	if err := u.userRepository.Exists(c, map[string]any{"phone_number": userModel.PhoneNumber}); err != nil {
		return model.User{}, err
	}

	newUser, err := u.userRepository.Add(c, userModel)
	if err != nil {
		return model.User{}, err
	}

	if _, err = u.emailService.AddEmail(c, newUser.ID, userModel.Email); err != nil {
		return model.User{}, err
	}

	if err := u.emailService.AproveEmail(c, newUser.ID, userModel.Email); err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (u UserService) UpdateUser(c context.Context, userModel model.User) (model.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserService[UpdateUser]")
	defer span.Finish()

	userRecord, err := u.GetUser(c, map[string]any{"id": userModel.ID})

	if err != nil {
		return model.User{}, err
	}

	if userRecord.Email != userModel.Email && len(userModel.Email) > 0 {
		if err := u.emailService.EmailExists(c, userModel.Email); err != nil {
			return model.User{}, err
		}
	}

	if userRecord.NationalId != userModel.NationalId && len(userModel.NationalId) > 0 {
		if err := u.userRepository.Exists(c, map[string]any{"national_id": userModel.NationalId}); err != nil {
			return model.User{}, err
		}
	}

	if userRecord.PhoneNumber != userModel.PhoneNumber && len(userModel.PhoneNumber) > 0 {
		if err := u.userRepository.Exists(c, map[string]any{"phone_number": userModel.PhoneNumber}); err != nil {
			return model.User{}, err
		}
	}

	_, err = u.userRepository.Update(c, userModel)
	if err != nil {
		return model.User{}, err
	}

	userRecord, err = u.GetUser(c, map[string]any{"id": userModel.ID})
	if err != nil {
		return model.User{}, err
	}

	return userRecord, nil
}

func (u UserService) DeleteUser(c context.Context, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "UserService[DeleteUser]")
	defer span.Finish()

	if _, err := u.GetUser(c, map[string]any{"id": userId}); err != nil {
		return err
	}

	return u.userRepository.Delete(c, userId)
}
