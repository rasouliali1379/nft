package user

import (
	"context"
	"errors"
	contract "nft/contract"
	merror "nft/error"
	"nft/infra/jtrace"
	"nft/infra/persist/model"
	userentity "nft/internal/user/entity"
	usermodel "nft/internal/user/model"
	"nft/pkg/crypt"
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

func (u UserRepository) Exists(c context.Context, conditions persist.Conds) (bool, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[Exists]")
	defer span.Finish()

	if _, err := u.db.Get(c, &userentity.User{}, conditions); err != nil {
		if errors.Is(err, merror.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (u UserRepository) Add(c context.Context, model usermodel.User) (usermodel.User, error) {
	span, ctx := jtrace.T().SpanFromContext(c, "UserRepository[Add]")
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
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[Update]")
	defer span.Finish()

	data := createMapFromUserModel(userModel)

	updatedUser, err := u.db.Update(c, &userentity.User{ID: userModel.ID}, data)
	if err != nil {
		return usermodel.User{}, err
	}

	return mapUserEntityToModel(updatedUser.(*userentity.User)), nil
}

func (u UserRepository) Delete(c context.Context, userId uuid.UUID) error {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[Delete]")
	defer span.Finish()

	if _, err := u.db.Update(c, &userentity.User{ID: userId}, map[string]any{"deleted_at": time.Now()}); err != nil {
		return err
	}
	return nil
}

func (u UserRepository) Get(c context.Context, conditions persist.Conds) (usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[Get]")
	defer span.Finish()

	user, err := u.db.Get(c, &userentity.User{}, conditions)
	if err != nil {
		return usermodel.User{}, err
	}
	return mapUserEntityToModel(user.(*userentity.User)), nil
}

func (u UserRepository) GetAll(c context.Context) ([]usermodel.User, error) {
	span, c := jtrace.T().SpanFromContext(c, "UserRepository[GetAll]")
	defer span.Finish()

	userList, err := u.db.GetAll(c, &[]userentity.User{}, map[string]any{})

	if err != nil {
		return nil, err
	}

	return createUserModelList(userList.(*[]userentity.User)), nil
}
