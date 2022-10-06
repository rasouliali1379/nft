package app

import (
	"context"
	"fmt"
	"log"
	"nft/client/jtrace"
	"nft/client/persist"
	"nft/client/server"
	"nft/client/storage"
	"nft/config"
	"nft/contract"
	"nft/src/talan"
	"os"
	"syscall"
	"time"

	"go.uber.org/fx"

	//modules
	"nft/src/auth"
	"nft/src/card"
	"nft/src/category"
	"nft/src/collection"
	"nft/src/email"
	"nft/src/file"
	"nft/src/jwt"
	"nft/src/kyc"
	"nft/src/nft"
	"nft/src/otp"
	"nft/src/user"
)

func Start() {
	fmt.Println("\n\n--------------------------------")
	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	for {
		fxNew := fx.New(
			fx.Provide(server.New),
			fx.Provide(persist.New),
			fx.Provide(storage.New),

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
