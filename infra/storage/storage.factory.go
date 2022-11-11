package storage

import (
	"context"
	"log"
	"nft/contract"
	"nft/infra/storage/aws"

	"go.uber.org/fx"
)

func New(lc fx.Lifecycle) contract.IStorage {
	log.Println("initialing minio storage")
	var storage aws.Aws
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {

			if err := storage.Init(c); err != nil {
				return err
			}
			log.Println("minio storage initialized successfully")
			return nil
		},
	})
	return &storage
}
