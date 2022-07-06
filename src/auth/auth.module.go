package auth

import (
	"go.uber.org/fx"
	contract "maskan/contract"
)

var Module = fx.Options(
	fx.Provide(NewAuthRepository),
	fx.Provide(func(params AuthServiceParams) contract.IAuthService {
		return NewAuthService(params)
	}),
	fx.Provide(NewAuthController),
)