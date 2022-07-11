package app

import (
	"context"
	"fmt"
	"log"
	"maskan/app/server"

	// "maskan/client/elk"
	// "maskan/client/jtrace"
	"maskan/client/jtrace"
	"maskan/client/persist"
	// "maskan/pkg/logger"
	"maskan/src/auth"
	"maskan/src/email"
	"maskan/src/jwt"
	"maskan/src/otp"
	"maskan/src/user"
	"os"
	"time"

	"maskan/config"

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
