package kyc

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewKycController),
	fx.Provide(NewKYCService),
	fx.Provide(NewKYCRepository),
)
