package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"nft/client/persist"
	"nft/client/server"
	"nft/client/storage"
	"nft/config"
	"nft/contract"
	"nft/src/auth"
	authdto "nft/src/auth/dto"
	"nft/src/card"
	"nft/src/category"
	"nft/src/collection"
	"nft/src/email"
	"nft/src/file"
	"nft/src/jwt"
	jwtmodel "nft/src/jwt/model"
	"nft/src/kyc"
	"nft/src/nft"
	"nft/src/otp"
	"nft/src/user"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var token string

var _ = BeforeSuite(func() {
	err := fx.New(
		fx.Provide(persist.New),
		fx.Provide(storage.New),
		fx.Provide(server.New),

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

		fx.Invoke(config.InitConfigs),
		fx.Invoke(migrate),
		fx.Invoke(serve),
	).Start(context.Background())
	if err != nil {
		return
	}

	client := resty.New()
	baseUrl := fmt.Sprintf("http://%s:%s/v1/auth/", config.C().App.Http.Host, config.C().App.Http.Port)

	signUpDto := authdto.SignUpRequest{
		FirstName:      "Ali",
		LastName:       "Rasouli",
		NationalId:     "324242344324",
		Email:          "test@gmail.com",
		PhoneNumber:    "234234234324",
		LandLineNumber: "02133073333",
		Password:       "ali1379",
		Province:       "tehran",
		City:           "terhan",
		Address:        "mahallati",
	}

	resp, err := client.R().
		SetBody(signUpDto).
		Post(baseUrl + "signup")

	if err != nil {
		AbortSuite(fmt.Sprintf("failed to verify user email for testing categories: %s", err.Error()))
	}

	var signUpResponse authdto.OtpToken
	err = json.Unmarshal(resp.Body(), &signUpResponse)
	if err != nil {
		AbortSuite(fmt.Sprintf("failed to unmarshal sign up response for testing categories: %s", err.Error()))
	}

	resp, err = client.R().
		SetBody(authdto.VerifyEmailRequest{
			Token: signUpResponse.Token,
			Code:  "111111",
		}).
		Post(baseUrl + "verify-email")

	if err != nil {
		AbortSuite(fmt.Sprintf("failed to verify user email for testing categories: %s", err.Error()))
	}

	loginDto := authdto.LoginRequest{
		Email:    "test@gmail.com",
		Password: "ali1379",
	}

	resp, err = client.R().
		SetBody(loginDto).
		Post(baseUrl + "login")

	if err != nil {
		AbortSuite(fmt.Sprintf("failed to unmarshal login response: %s", err.Error()))
	}

	var jwtToken jwtmodel.Jwt
	err = json.Unmarshal(resp.Body(), &jwtToken)
	if err != nil {
		AbortSuite(fmt.Sprintf("failed unmarshal jwt struct: %s", err.Error()))
	}
	token = jwtToken.AccessToken
})

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
