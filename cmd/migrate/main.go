package main

import (
	"context"
	"nft/config"
	"nft/contract"
	"nft/infra/persist"

	"go.uber.org/fx"
)

func main() {
	if err := fx.New(
		fx.Provide(persist.New),
		fx.Invoke(config.InitConfigs),
		fx.Invoke(migrate),
	).Start(context.TODO()); err != nil {
		panic(err)
	}
}

func migrate(lc fx.Lifecycle, db contract.IPersist) {
	lc.Append(
		fx.Hook{
			OnStart: func(c context.Context) error {
				return db.Migrate(c)
			},
		},
	)
}
