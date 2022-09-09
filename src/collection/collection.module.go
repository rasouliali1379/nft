package collection

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewCollectionController),
	fx.Provide(NewCollectionService),
	fx.Provide(NewCollectionRepository),
)
