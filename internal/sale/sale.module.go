package sale

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewSaleController),
	fx.Provide(NewSaleService),
	fx.Provide(NewSaleRepository),
)
