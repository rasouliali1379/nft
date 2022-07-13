package user

import (
	"context"
	"errors"

	"nft/client/jtrace"
	contract "nft/contract"
	merror "nft/error"
	"nft/pkg/crypt"
	userentity "nft/src/user/entity"
	usermodel "nft/src/user/model"
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

func (u UserRepository) Exists(c context.Context, conditions map[string]any) error {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[UserExists]")
	defer span.Finish()

	if err := u.db.Exists(c, &userentity.User{}, conditions); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (u UserRepository) Add(c context.Context, model usermodel.User) (usermodel.User, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "UserRepository[AddUser]")
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

func (u UserRepository) Update(c context.Context, userModel usermodel.User) (usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[UpdateUser]")
	defer span.Finish()

	data := createMapFromUserModel(userModel)

	updatedUser, err := u.db.Update(c, &userentity.User{ID: userModel.ID}, data)
	if err != nil {
		return usermodel.User{}, err
	}

	return mapUserEntityToModel(updatedUser.(*userentity.User)), nil
}

func (u UserRepository) Delete(c context.Context, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[DeleteUser]")
	defer span.Finish()

	if _, err := u.db.Update(c, &userentity.User{ID: userId}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return err
	}
	return nil
}

func (u UserRepository) Get(c context.Context, conditions map[string]any) (usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[GetUser]")
	defer span.Finish()

	user, err := u.db.Get(c, &userentity.User{}, conditions)
	if err != nil {
		return usermodel.User{}, err
	}
	return mapUserEntityToModel(user.(*userentity.User)), nil
}

func (u UserRepository) GetAll(c context.Context) ([]usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[GetAllUsers]")
	defer span.Finish()

	userList, err := u.db.Get(c, &[]userentity.User{}, map[string]any{})

	if err != nil {
		return nil, err
	}

	return createUserModelList(userList.(*[]userentity.User)), nil
}
