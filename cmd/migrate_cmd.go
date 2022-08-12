package cmd

import (
	"context"
	"nft/client/persist"
	"nft/config"
	"nft/contract"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var migrateCMD = cobra.Command{
	Use:     "migrate",
	Long:    "migrate database strucutures. This will migrate tables",
	Aliases: []string{"m"},
	Run:     Runner.Migrate,
}

func (c *command) Migrate(cmd *cobra.Command, args []string) {
	fx.New(
		fx.Provide(persist.New),
		fx.Invoke(config.InitConfigs),
		// fx.Invoke(logger.InitGlobalLogger),
		fx.Invoke(migrate),
	).Start(context.TODO())
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
