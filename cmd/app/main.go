package main

import (
	"context"
	"fmt"
	"log"
	"nft/config"
	"nft/contract"
	"nft/infra/jtrace"
	"nft/infra/persist"
	"nft/infra/server"
	"nft/infra/storage"
	"nft/internal/offer"
	"nft/internal/sale"
	"nft/internal/talan"
	"nft/internal/transaction"
	"os"
	"syscall"
	"time"

	"go.uber.org/fx"

	//modules
	"nft/internal/auth"
	"nft/internal/card"
	"nft/internal/category"
	"nft/internal/collection"
	"nft/internal/email"
	"nft/internal/file"
	"nft/internal/jwt"
	"nft/internal/kyc"
	"nft/internal/nft"
	"nft/internal/otp"
	"nft/internal/user"
)

func main() {
	fmt.Println("\n\n--------------------------------")
	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	for {
		fxNew := fx.New(
			fx.Provide(server.New),
			fx.Provide(persist.New),
			fx.Provide(storage.New),

			sale.Module,
			auth.Module,
			user.Module,
			jwt.Module,
			otp.Module,
			email.Module,
			collection.Module,
			category.Module,
			kyc.Module,
			card.Module,
			nft.Module,
			file.Module,
			talan.Module,
			offer.Module,
			transaction.Module,

			fx.Invoke(initConfig),
			fx.Invoke(jtrace.InitGlobalTracer),
			fx.Invoke(migrate),
			fx.Invoke(serve),
		)

		startCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := fxNew.Start(startCtx); err != nil {
			log.Println(err)
			break
		}

		if val := <-fxNew.Done(); val == syscall.SIGTERM {
			break
		}

		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := fxNew.Stop(stopCtx); err != nil {
			log.Println(err)
			break
		}
	}
}

func initConfig(down fx.Shutdowner) {
	path, err := os.Getwd()
	if err != nil {
		panic("unable to initialize config")
	}
	config.InitConfigs(down, path)
}

func serve(lc fx.Lifecycle, server contract.IServer) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return server.ListenAndServe()
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown()
		},
	})
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
