package transaction

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewTransactionService),
	fx.Provide(NewTransactionRepository),
)
