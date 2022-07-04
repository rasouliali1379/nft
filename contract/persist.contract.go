package contract

import (
	"context"
	jwt "maskan/src/jwt/entity"
	user "maskan/src/user/entity"
)

type IPersist interface {
	Init(c context.Context) error
	Migrate(c context.Context) error
	Close(c context.Context) error

	UserExists(c context.Context, columnName string, value string) error
	GetAllUsers(c context.Context, conditions map[string]any) ([]user.User, error)
	GetUser(c context.Context, columnName string, value string) (user.User, error)
	CreateUser(c context.Context, entity user.User) (user.User, error)
	UpdateUser(c context.Context, entity user.User) (user.User, error)
	DeleteUser(c context.Context, userId string) error

	SaveToken(c context.Context, entity jwt.Jwt) error
	UpdateToken(c context.Context, id uint, token string) error
	RetrieveToken(c context.Context, token string) (jwt.Jwt, error)
}
