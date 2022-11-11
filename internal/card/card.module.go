package card

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewCardRepository),
	fx.Provide(NewCardService),
	fx.Provide(NewCardController),
)
