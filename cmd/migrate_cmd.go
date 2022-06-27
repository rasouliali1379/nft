package cmd

import (
	"context"
	"log"
	"maskan/client/persist"
	"maskan/config"

	user "maskan/src/user"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var migrateCMD = cobra.Command{
	Use:     "migrate",
	Long:    "migrate database strucutures. This will migrate tables",
	Aliases: []string{"m"},
	Run:     Runner.Migrate,
}

// migrate database with fake data
func (c *command) Migrate(cmd *cobra.Command, args []string) {
	fx.New(
		fx.Provide(persist.New),
		fx.Invoke(config.InitConfigs),
		// fx.Invoke(logger.InitGlobalLogger),
		fx.Invoke(migrateReducer(args)),
	).Start(context.TODO())
}

func migrateReducer(args []string) func(fx.Lifecycle, *gorm.DB) {
	return func(l fx.Lifecycle, d *gorm.DB) {
		migrate(l, d, args)
	}
}

func migrate(lc fx.Lifecycle, db *gorm.DB, _ []string) {
	lc.Append(fx.Hook{OnStart: func(_ context.Context) error {
		db.Transaction(func(tx *gorm.DB) error {
			log.Println("starting migration")
			if err := tx.AutoMigrate(&user.User{}); err != nil {
				log.Printf("migration failed with error: %s\n", err.Error())
				return err
			}

			log.Println("migration done successfully")
			return nil
		})
		return nil
	}})
}
