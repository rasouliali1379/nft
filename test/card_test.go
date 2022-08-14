package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nft/config"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	carddto "nft/src/card/dto"
)

var _ = Describe("Card Management", Ordered, func() {

	var baseUrl string
	card := carddto.AddCardRequest{
		CardNumber: "31230213809",
		IBAN:       "2131123122",
	}

	var cardId uuid.UUID

	client := resty.New()

	BeforeAll(func() {
		baseUrl = fmt.Sprintf("http://%s:%s/v1/card/", config.C().App.Http.Host, config.C().App.Http.Port)
	})

	Describe("add new card", func() {
		It("should add new card successfully", func() {

			resp, err := client.R().
				SetBody(card).
				SetAuthToken(token).
				Post(baseUrl)
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
		})
	})

	Describe("Get cards list", func() {
		It("should get cards list successfully", func() {
			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl)
			if err != nil {
				Fail("unable to make request to get card list", 3)
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var cardList carddto.CardList
			if err = json.Unmarshal(resp.Body(), &cardList); err != nil {
				Fail("unable to unmarshal card list", 3)
			}

			By("cards list should have one item")
			Expect(len(cardList.Cards)).To(Equal(1))

			cardId = cardList.Cards[0].ID
		})
	})

	Describe("Approve card", func() {
		It("should approve card successfully", func() {
			resp, err := client.R().
				SetAuthToken(token).
				Post(baseUrl + cardId.String() + "/approve")
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Get card", func() {
		It("should get single card successfully", func() {

			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl + cardId.String())
			Expect(err).NotTo(HaveOccurred())

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
				SetAuthToken(token).
				Delete(baseUrl + cardId.String())
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
