package kyc

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewKYCController),
	fx.Provide(NewKYCService),
	fx.Provide(NewKYCRepository),
)
