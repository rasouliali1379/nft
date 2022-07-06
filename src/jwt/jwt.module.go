package jwt

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewJwtRepository),
	fx.Provide(NewJwtService),
)
