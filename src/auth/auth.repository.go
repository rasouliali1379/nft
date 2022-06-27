package auth

import (
	"maskan/client/persist"
	contract "maskan/src/auth/contract"

	"go.uber.org/fx"
)

type AuthRepository struct {
	db persist.IPersist
}

type AuthRepositoryParams struct {
	fx.In
	DB persist.IPersist
}

func NewAuthRepository(params AuthRepositoryParams) contract.IAuthRepository {
	return &AuthRepository{
		db: params.DB,
	}
}