package persist

import (
	"context"
	"log"
	"maskan/client/persist/postgres"

	"go.uber.org/fx"
)

type IPersist interface {
	Init() error
	Close() error
}

func New(lc fx.Lifecycle) IPersist {

	var db IPersist
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			db = postgres.Postgres{}

			if err := db.Init(); err != nil {
				return err
			}
			log.Println("postgres database loaded successfully")
			return nil
		},
		OnStop: func(_ context.Context) error {
			if err := db.Close(); err != nil {
				return err
			}
			log.Println("postgres database connection closed")
			return nil
		},
	})
	return db
}
