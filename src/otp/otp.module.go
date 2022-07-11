package otp

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewOtpService),
	fx.Provide(NewOtpRepository),
)
