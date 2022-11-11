package offer

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewOfferController),
	fx.Provide(NewOfferService),
	fx.Provide(NewOfferRepository),
)
