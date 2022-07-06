package user

import (
	"go.uber.org/fx"
	"maskan/contract"
)

var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(func(params UserServiceParams) contract.IUserService {
		return NewUserService(params)
	}),
	fx.Provide(NewUserController),
)
