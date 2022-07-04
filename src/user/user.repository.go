package user

import (
	"context"
	"errors"
	"fmt"
	"maskan/client/jtrace"
	contract "maskan/contract"
	merror "maskan/error"
	"maskan/pkg/crypt"
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

func (u UserRepository) AddUser(c context.Context, model usermodel.User) (usermodel.User, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "repository[AddUser]")
	defer span.Finish()

	mappedEntity := mapUserModelToEntity(model)

	password, err := crypt.Hash(model.Password)
	if err != nil {
		return usermodel.User{}, err
	}

	mappedEntity.ID = uuid.New()
	mappedEntity.Password = password

	userEntity, err := u.db.CreateUser(ctx, mappedEntity)
	if err != nil {
		return usermodel.User{}, err
	}

	return mapUserEntityToModel(userEntity), nil
}

func (u UserRepository) UpdateUser(c context.Context, userModel usermodel.User) (usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[UpdateUser]")
	defer span.Finish()

	updatedUser, err := u.db.UpdateUser(c, mapUserModelToEntity(userModel))
	if err != nil {
		return usermodel.User{}, fmt.Errorf("error happened while updating jwt: %w", err)
	}

	return mapUserEntityToModel(updatedUser), nil
}

func (u UserRepository) DeleteUser(c context.Context, userId string) error {
	span, c := jtrace.T().SpanFromContext(c, "repository[DeleteUser]")
	defer span.Finish()

	if err := u.db.DeleteUser(c, userId); err != nil {
		return fmt.Errorf("error happened while deleting a user: %w", err)
	}
	return nil
}

func (u UserRepository) GetUser(c context.Context, query usermodel.UserQuery) (usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[GetUserByEmail]")
	defer span.Finish()
	if len(query.ID) > 0 {
		user, err := u.db.GetUser(c, "id", query.ID)
		if err != nil {
			return usermodel.User{}, err
		}
		return mapUserEntityToModel(user), nil
	} else if len(query.Email) > 0 {
		user, err := u.db.GetUser(c, "email", query.Email)
		if err != nil {
			return usermodel.User{}, err
		}
		return mapUserEntityToModel(user), nil
	} else if len(query.NationalId) > 0 {
		user, err := u.db.GetUser(c, "national_id", query.NationalId)
		if err != nil {
			return usermodel.User{}, err
		}
		return mapUserEntityToModel(user), nil
	}

	return usermodel.User{}, merror.ErrNoQueries
}

func (u UserRepository) GetAllUsers(c context.Context) ([]usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[GetAllUsers]")
	defer span.Finish()

	userList, err := u.db.GetAllUsers(c, map[string]any{})

	if err != nil {
		return nil, err
	}

	return createUserModelList(userList), nil
}
