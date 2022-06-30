package user

import (
	"context"
	"errors"
	"fmt"
	"maskan/client/jtrace"
	contract "maskan/contract"
	merror "maskan/error"
	"maskan/pkg/crypt"
	model "maskan/src/auth/model"
	usermodel "maskan/src/user/model"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type UserRepository struct {
	db contract.IPersist
}

type UserRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewUserRepository(params UserRepositoryParams) contract.IUserRepository {
	return &UserRepository{
		db: params.DB,
	}
}

func (u UserRepository) NationalIdExists(c context.Context, nationalId string) error {
	span, ctx := jtrace.T().SpanFromContext(c, "repository[NationalIdExists]")
	defer span.Finish()

	if err := u.db.UserExists(ctx, "national_id", nationalId); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("error happened while searching for a national id: %w", err)
	}

	return merror.ErrNationalIdExists
}

func (u UserRepository) PhoneNumberExists(c context.Context, phoneNumber string) error {
	span, ctx := jtrace.T().SpanFromContext(c, "repository[PhoneNumberExists]")
	defer span.Finish()

	if err := u.db.UserExists(ctx, "phone_number", phoneNumber); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("error happened while searching for a phone number: %w", err)
	}

	return merror.ErrPhoneNumberExists
}

func (u UserRepository) EmailExists(c context.Context, email string) error {
	span, ctx := jtrace.T().SpanFromContext(c, "repository[EmailExists]")
	defer span.Finish()

	if err := u.db.UserExists(ctx, "email", email); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("error happened while searching for an email: %w", err)
	}

	return merror.ErrEmailExists
}

func (u UserRepository) AddUser(c context.Context, dto model.SignUpRequest) (string, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "repository[AddUser]")
	defer span.Finish()

	user := mapSignUpRequestModelToEntity(dto)

	password, err := crypt.Hash(dto.Password)
	if err != nil {
		return "", err
	}

	user.ID = uuid.New()
	user.Password = password

	userId, err := u.db.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (u UserRepository) GetUser(c context.Context, email string) (usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[GetUser]")
	defer span.Finish()

	user, err := u.db.GetUser(c, "email" , email)
	if err != nil {
		return usermodel.User{}, err
	}

	return mapUserEntityToModel(user), nil
}
