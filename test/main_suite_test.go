package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"nft/config"
	"nft/contract"
	"nft/infra/persist"
	"nft/infra/server"
	"nft/infra/storage"
	"nft/internal/auth"
	authdto "nft/internal/auth/dto"
	"nft/internal/card"
	"nft/internal/category"
	"nft/internal/collection"
	"nft/internal/email"
	"nft/internal/file"
	"nft/internal/jwt"
	jwtmodel "nft/internal/jwt/model"
	"nft/internal/kyc"
	"nft/internal/nft"
	"nft/internal/offer"
	"nft/internal/otp"
	"nft/internal/sale"
	"nft/internal/talan"
	"nft/internal/transaction"
	"nft/internal/user"
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
		talan.Module,
		offer.Module,
		sale.Module,
		transaction.Module,

		fx.Invoke(initConfig),
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
	log.Println(token)
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

func initConfig(down fx.Shutdowner) {
	config.InitConfigs(down, ".")
}
