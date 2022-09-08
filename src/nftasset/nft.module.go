package nftasset

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewNftController),
	fx.Provide(NewNftService),
	fx.Provide(NewNftRepository),
)
