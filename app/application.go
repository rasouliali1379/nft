package app

import (
	"context"
	"fmt"
	"log"
	"nft/client/jtrace"
	"nft/client/persist"
	"nft/client/server"
	"nft/config"
	"nft/contract"
	"os"
	"time"

	//modules
	"nft/src/auth"
	"nft/src/category"
	"nft/src/collection"
	"nft/src/email"
	"nft/src/jwt"
	"nft/src/otp"
	"nft/src/user"

	"go.uber.org/fx"
)

func Start() {
	fmt.Println("\n\n--------------------------------")
	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	for {
		fxNew := fx.New(
			fx.Provide(persist.New),
			auth.Module,
			user.Module,
			jwt.Module,
			otp.Module,
			email.Module,
			collection.Module,
			category.Module,
			fx.Provide(server.New),
			fx.Invoke(config.InitConfigs),
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
		if val := <-fxNew.Done(); val == os.Interrupt {
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
