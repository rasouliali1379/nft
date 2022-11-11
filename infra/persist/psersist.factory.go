package persist

import (
	"context"
	"log"
	"nft/contract"
	"nft/infra/persist/postgres"

	"go.uber.org/fx"
)

func New(lc fx.Lifecycle) contract.IPersist {

	var db postgres.Postgres
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {

			if err := db.Init(c); err != nil {
				return err
			}
			log.Println("postgres database loaded successfully")
			return nil
		},
		OnStop: func(c context.Context) error {
			if err := db.Close(c); err != nil {
				return err
			}
			log.Println("postgres database connection closed")
			return nil
		},
	})
	return &db
}
