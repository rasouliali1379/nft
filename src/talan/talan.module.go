package talan

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewTalanRepository),
	fx.Provide(NewTalanService),
)
