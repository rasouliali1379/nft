package user

import (
	"nft/contract"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(func(params UserServiceParams) contract.IUserService {
		return NewUserService(params)
	}),
	fx.Provide(NewUserController),
)
