package test

import (
	"context"
	"maskan/app/server"
	"maskan/client/persist"
	"maskan/config"
	"maskan/contract"
	"maskan/src/auth"
	"maskan/src/jwt"
	"maskan/src/otp"
	"maskan/src/user"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

func TestTest(t *testing.T) {
	fx.New(
		fx.Provide(persist.New),
		auth.Module,
		user.Module,
		jwt.Module,
		otp.Module,
		fx.Provide(server.New),
		fx.Invoke(config.InitConfigs),
		fx.Invoke(migrate),
		fx.Invoke(serve),
	).Start(context.TODO())

	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
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
