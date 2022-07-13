package test

import (
	"context"
	"nft/app/server"
	"nft/client/persist"
	"nft/config"
	"nft/contract"
	"nft/src/auth"
	"nft/src/email"
	"nft/src/jwt"
	"nft/src/otp"
	"nft/src/user"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestTest(t *testing.T) {
	fxtest.New(
		t,
		fx.Provide(persist.New),
		auth.Module,
		user.Module,
		jwt.Module,
		otp.Module,
		email.Module,
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
