package email

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewEmailRepository),
	fx.Provide(NewEmailService),
)
