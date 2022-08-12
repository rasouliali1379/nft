package file

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewFileRepository),
	fx.Provide(NewFileService),
)
