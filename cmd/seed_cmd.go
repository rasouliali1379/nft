package cmd

import (
	"context"
	"log"
	"nft/client/persist"
	"nft/config"
	"nft/pkg/logger"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	seedCMD = cobra.Command{
		Use:  "seed database",
		Long: "seed database strucutures. This will seed tables",
		Run:  Runner.Seed,
	}
)

// seed database with fake data
func (c *command) Seed(cmd *cobra.Command, args []string) {
	fx.New(
		fx.Provide(persist.New),
		fx.Invoke(config.InitConfigs),
		fx.Invoke(logger.InitGlobalLogger),
		fx.Invoke(seed),
	).Start(context.TODO())
}

func seed(lc fx.Lifecycle, db *pg.DB) {
	// Do all your seeding here
	lc.Append(fx.Hook{OnStart: func(c context.Context) error {
		log.Println("Data seeded successfully")
		return nil
	}})
}
