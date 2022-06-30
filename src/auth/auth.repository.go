package auth

import (
	contract "maskan/contract"

	"go.uber.org/fx"
)

type AuthRepository struct {
	db contract.IPersist
}

type AuthRepositoryParams struct {
	fx.In
	DB contract.IPersist
}

func NewAuthRepository(params AuthRepositoryParams) contract.IAuthRepository {
	return &AuthRepository{
		db: params.DB,
	}
}
