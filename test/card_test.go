package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nft/config"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	authdto "nft/src/auth/dto"
	carddto "nft/src/card/dto"
	jwt "nft/src/jwt/model"
)

var cardToken string

var _ = Describe("Card Management", Ordered, func() {

	BeforeAll(func() {
		client := resty.New()
		baseUrl := fmt.Sprintf("http://%s:%s/v1/auth/", config.C().App.Http.Host, config.C().App.Http.Port)

		signUpDto := authdto.SignUpRequest{
			FirstName:      "Ali",
			LastName:       "Rasouli",
			NationalId:     "0123456782",
			Email:          "testcard@gmail.com",
			PhoneNumber:    "09368045734",
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
			Fail(fmt.Sprintf("failed to verify user email for testing cardegories: %s", err.Error()), 5)
		}

		var signUpResponse authdto.OtpToken
		err = json.Unmarshal(resp.Body(), &signUpResponse)
		if err != nil {
			Fail(fmt.Sprintf("failed to unmarshal sign up response for testing cardegories: %s", err.Error()), 5)
		}

		resp, err = client.R().
			SetBody(authdto.VerifyEmailRequest{
				Token: signUpResponse.Token,
				Code:  "111111",
			}).
			Post(baseUrl + "verify-email")

		if err != nil {
			Fail(fmt.Sprintf("failed to verify user email for testing cardegories: %s", err.Error()), 5)
		}

		loginDto := authdto.LoginRequest{
			Email:    "testcard@gmail.com",
			Password: "ali1379",
		}

		resp, err = client.R().
			SetBody(loginDto).
			Post(baseUrl + "login")

		if err != nil {
			Fail(fmt.Sprintf("failed to unmarshal login response: %s", err.Error()), 5)
		}
		log.Println(resp)
		var jwtToken jwt.Jwt
		err = json.Unmarshal(resp.Body(), &jwtToken)
		if err != nil {
			Fail(fmt.Sprintf("failed unmarshal jwt struct: %s", err.Error()), 5)
		}

		cardToken = jwtToken.AccessToken
	})

	card := carddto.AddCardRequest{
		CardNumber: "31230213809",
		IBAN:       "2131123122",
	}

	var cardId uuid.UUID

	client := resty.New()
	baseUrl := fmt.Sprintf("http://%s:%s/v1/card/", config.C().App.Http.Host, config.C().App.Http.Port)

	Describe("add new card", func() {
		It("should add new card successfully", func() {
			resp, err := client.R().
				SetBody(card).
				SetAuthToken(cardToken).
				Post(baseUrl)

			Expect(err).NotTo(HaveOccurred())

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
		})
	})

	Describe("Get cards list", func() {
		It("should get cards list successfully", func() {
			resp, err := client.R().
				SetAuthToken(cardToken).
				Get(baseUrl)

			if err != nil {
				Fail(fmt.Sprintf("failed to make request to get cards list: %s", err.Error()))
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var cardList carddto.CardList
			if err = json.Unmarshal(resp.Body(), &cardList); err != nil {
				Fail("unable to unmarshal card list")
			}

			By("cards list should have one item")
			Expect(len(cardList.Cards)).To(Equal(1))

			cardId = cardList.Cards[0].ID
		})
	})

	Describe("Approve card", func() {
		It("should approve card successfully", func() {
			resp, err := client.R().
				SetAuthToken(cardToken).
				Post(baseUrl + cardId.String() + "/approve")

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Get card", func() {
		It("should get single card successfully", func() {

			if cardId == uuid.Nil {
				Fail("card id is empty")
			}

			resp, err := client.R().
				SetAuthToken(cardToken).
				Get(baseUrl + cardId.String())

			if err != nil {
				Fail(fmt.Sprintf("failed to make the request to get card: %s", err.Error()))
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var card carddto.Card
			if err = json.Unmarshal(resp.Body(), &card); err != nil {
				Fail("unable to unmarshal card object")
			}

			By("approved field should be true")
			Expect(card.Approved).To(BeTrue())
		})
	})

	Describe("Delete card", func() {
		It("should delete user successfully", func() {

			resp, err := client.R().
				SetAuthToken(cardToken).
				Delete(baseUrl + cardId.String())

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
