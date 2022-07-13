package app

import (
	"context"
	"fmt"
	"log"
	"nft/app/server"
	"nft/contract"

	// "nft/client/elk"
	// "nft/client/jtrace"
	"nft/client/jtrace"
	"nft/client/persist"

	// "nft/pkg/logger"
	"nft/src/auth"
	"nft/src/email"
	"nft/src/jwt"
	"nft/src/otp"
	"nft/src/user"
	"os"
	"time"

	"nft/config"

	"go.uber.org/fx"
)

// StartApplication func
func Start() {
	fmt.Println("\n\n--------------------------------")
	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init configs

	for {
		fxNew := fx.New(
			//		fx.Provide(broker.NewNats),
			//		fx.Provide(redis.NewRedis),
			// fx.Provide(elk.NewLogStash),
			fx.Provide(persist.New),
			auth.Module,
			user.Module,
			jwt.Module,
			otp.Module,
			email.Module,
			fx.Provide(server.New),
			fx.Invoke(config.InitConfigs),
			// fx.Invoke(logger.InitGlobalLogger),
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

func serve(lc fx.Lifecycle, server server.IServer) {
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
