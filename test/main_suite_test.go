package test

import (
	"context"
	"nft/client/persist"
	"nft/client/server"
	"nft/config"
	"nft/contract"
	"nft/src/auth"
	"nft/src/category"
	"nft/src/collection"
	"nft/src/email"
	"nft/src/jwt"
	"nft/src/otp"
	"nft/src/user"
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
		email.Module,
		collection.Module,
		category.Module,
		fx.Provide(server.New),
		fx.Invoke(config.InitConfigs),
		fx.Invoke(migrate),
		fx.Invoke(serve),
	).Start(context.Background())
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
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
