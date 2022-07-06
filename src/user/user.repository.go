package user

import (
	"context"
	"errors"
	"fmt"
	"maskan/client/jtrace"
	contract "maskan/contract"
	merror "maskan/error"
	"maskan/pkg/crypt"
	userentity "maskan/src/user/entity"
	usermodel "maskan/src/user/model"
	"time"

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

func (u UserRepository) UserExists(c context.Context, query usermodel.UserQuery) error {
	span, c := jtrace.T().SpanFromContext(c, "repository[UserExists]")
	defer span.Finish()

	if len(query.Email) > 0 {
		if err := u.db.Exists(c, &userentity.User{}, map[string]any{"email": query.Email}); err != nil {
			if errors.Is(err, merror.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("error happened while searching for an email: %w", err)
		}
		return nil
	} else if len(query.NationalId) > 0 {
		if err := u.db.Exists(c, &userentity.User{}, map[string]any{"national_id": query.NationalId}); err != nil {
			if errors.Is(err, merror.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("error happened while searching for a national id: %w", err)
		}
		return nil
	} else if len(query.PhoneNumber) > 0 {
		if err := u.db.Exists(c, &userentity.User{}, map[string]any{"phone_number": query.PhoneNumber}); err != nil {
			if errors.Is(err, merror.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("error happened while searching for a phone number: %w", err)
		}
		return nil
	}

	return merror.ErrNoQueries
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

	userEntity, err := u.db.Create(ctx, &mappedEntity)
	if err != nil {
		return usermodel.User{}, err
	}

	return mapUserEntityToModel(userEntity.(*userentity.User)), nil
}

func (u UserRepository) UpdateUser(c context.Context, userModel usermodel.User) (usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[UpdateUser]")
	defer span.Finish()

	data := createMapFromUserModel(userModel)

	updatedUser, err := u.db.Update(c, &userentity.User{ID: userModel.ID}, data)
	if err != nil {
		return usermodel.User{}, fmt.Errorf("error happened while updating jwt: %w", err)
	}

	return mapUserEntityToModel(updatedUser.(*userentity.User)), nil
}

func (u UserRepository) DeleteUser(c context.Context, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "repository[DeleteUser]")
	defer span.Finish()

	if _, err := u.db.Update(c, &userentity.User{ID: userId}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return fmt.Errorf("error happened while deleting a user: %w", err)
	}
	return nil
}

func (u UserRepository) GetUser(c context.Context, query usermodel.UserQuery) (usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[GetUser]")
	defer span.Finish()

	if query.ID != uuid.Nil {
		user, err := u.db.Get(c, &userentity.User{}, map[string]any{"id": query.ID})
		if err != nil {
			return usermodel.User{}, err
		}
		return mapUserEntityToModel(user.(*userentity.User)), nil
	} else if len(query.Email) > 0 {
		user, err := u.db.Get(c, &userentity.User{}, map[string]any{"email": query.Email})
		if err != nil {
			return usermodel.User{}, err
		}
		return mapUserEntityToModel(user.(*userentity.User)), nil
	} else if len(query.NationalId) > 0 {
		user, err := u.db.Get(c, &userentity.User{}, map[string]any{"national_id": query.NationalId})
		if err != nil {
			return usermodel.User{}, err
		}
		return mapUserEntityToModel(user.(*userentity.User)), nil
	}

	return usermodel.User{}, merror.ErrNoQueries
}

func (u UserRepository) GetAllUsers(c context.Context) ([]usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "repository[GetAllUsers]")
	defer span.Finish()

	userList, err := u.db.Get(c, &[]userentity.User{}, map[string]any{})

	if err != nil {
		return nil, err
	}

	return createUserModelList(userList.(*[]userentity.User)), nil
}
