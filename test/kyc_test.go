package test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nft/config"
	"os"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	auth "nft/src/auth/dto"
	jwt "nft/src/jwt/model"
)

var kycToken string

var _ = Describe("Kyc Management", Ordered, func() {
	var kycImage1, kycImage2 io.Reader
	BeforeAll(func() {
		client := resty.New()
		baseUrl := fmt.Sprintf("http://%s:%s/v1/auth/", config.C().App.Http.Host, config.C().App.Http.Port)

		signUpDto := auth.SignUpRequest{
			FirstName:      "Ali",
			LastName:       "Rasouli",
			NationalId:     "0123456785",
			Email:          "testkyc@gmail.com",
			PhoneNumber:    "09368045734",
			LandLineNumber: "02133073333",
			Password:       "ali1379",
			Province:       "tehran",
			City:           "Terhan",
			Address:        "mahallati",
		}

		resp, err := client.R().
			SetBody(signUpDto).
			Post(baseUrl + "signup")

		if err != nil {
			Fail(fmt.Sprintf("failed to verify user email for testing kycegories: %s", err.Error()), 5)
		}

		var signUpResponse auth.OtpToken
		err = json.Unmarshal(resp.Body(), &signUpResponse)
		if err != nil {
			Fail(fmt.Sprintf("failed to unmarshal sign up response for testing kycegories: %s", err.Error()), 5)
		}

		resp, err = client.R().
			SetBody(auth.VerifyEmailRequest{
				Token: signUpResponse.Token,
				Code:  "111111",
			}).
			Post(baseUrl + "verify-email")

		if err != nil {
			Fail(fmt.Sprintf("failed to verify user email for testing kycegories: %s", err.Error()), 5)
		}

		loginDto := auth.LoginRequest{
			Email:    "testkyc@gmail.com",
			Password: "ali1379",
		}

		resp, err = client.R().
			SetBody(loginDto).
			Post(baseUrl + "login")

		if err != nil {
			Fail(fmt.Sprintf("failed to unmarshal login response: %s", err.Error()), 5)
		}

		var jwtToken jwt.Jwt
		err = json.Unmarshal(resp.Body(), &jwtToken)
		if err != nil {
			Fail(fmt.Sprintf("failed unmarshal jwt struct: %s", err.Error()), 5)
		}
		kycToken = jwtToken.AccessToken

		path, err := os.Getwd()
		if err != nil {
			Fail(fmt.Sprintf("failed get current directory path: %s", err.Error()), 5)
		}

		kycImage1, err = os.Open(path + "/assets/kyc-test-image-1.jpeg")
		if err != nil {
			Fail(fmt.Sprintf("failed to open kyc image file 1: %s", err.Error()), 5)
		}

		kycImage2, err = os.Open(path + "/assets/kyc-test-image-2.jpg")
		if err != nil {
			Fail(fmt.Sprintf("failed to open kyc image file 2: %s", err.Error()), 5)
		}
	})

	//var kycId uuid.UUID

	client := resty.New()
	baseUrl := fmt.Sprintf("http://%s:%s/v1/kyc/", config.C().App.Http.Host, config.C().App.Http.Port)

	Describe("appeal for kyc", func() {
		It("should request for kyc successfully", func() {
			log.Println(kycToken)
			resp, err := client.R().
				SetMultipartFields(
					&resty.MultipartField{
						Param:  "id_card",
						Reader: kycImage1,
					},
					&resty.MultipartField{
						Param:  "portrait",
						Reader: kycImage2,
					}).
				SetAuthToken(kycToken).
				Post(baseUrl)

			if err != nil {
				AbortSuite(fmt.Sprintf("failed appeal for kyc: %s", err.Error()))
			}

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
		})
	})

	//Describe("Get kycs list", func() {
	//	It("should get kycs list successfully", func() {
	//		resp, err := client.R().
	//			SetAuthToken(kycToken).
	//			Get(baseUrl)
	//
	//		if err != nil {
	//			Fail(fmt.Sprintf("failed to make request to get kycs list: %s", err.Error()))
	//		}
	//
	//		By("status code should be 200")
	//		Expect(resp.StatusCode()).To(Equal(http.StatusOK))
	//
	//		var kycList kycdto.KycList
	//		if err = json.Unmarshal(resp.Body(), &kycList); err != nil {
	//			Fail("unable to unmarshal kyc list")
	//		}
	//
	//		By("kycs list should have one item")
	//		Expect(len(kycList.Kycs)).To(Equal(1))
	//
	//		kycId = kycList.Kycs[0].ID
	//	})
	//})
	//
	//Describe("Approve kyc", func() {
	//	It("should approve kyc successfully", func() {
	//		resp, err := client.R().
	//			SetAuthToken(kycToken).
	//			Post(baseUrl + kycId.String() + "/approve")
	//
	//		if err != nil {
	//			Expect(err).NotTo(HaveOccurred())
	//		}
	//
	//		By("status code should be 200")
	//		Expect(resp.StatusCode()).To(Equal(http.StatusOK))
	//	})
	//})
	//
	//Describe("Get kyc", func() {
	//	It("should get single kyc successfully", func() {
	//
	//		if kycId == uuid.Nil {
	//			Fail("kyc id is empty")
	//		}
	//
	//		resp, err := client.R().
	//			SetAuthToken(kycToken).
	//			Get(baseUrl + kycId.String())
	//
	//		if err != nil {
	//			Fail(fmt.Sprintf("failed to make the request to get kyc: %s", err.Error()))
	//		}
	//
	//		By("status code should be 200")
	//		Expect(resp.StatusCode()).To(Equal(http.StatusOK))
	//
	//		var kyc kycdto.Kyc
	//		if err = json.Unmarshal(resp.Body(), &kyc); err != nil {
	//			Fail("unable to unmarshal kyc object")
	//		}
	//
	//		By("approved field should be true")
	//		Expect(kyc.Approved).To(BeTrue())
	//	})
	//})
	//
	//Describe("Delete kyc", func() {
	//	It("should delete user successfully", func() {
	//
	//		resp, err := client.R().
	//			SetAuthToken(kycToken).
	//			Delete(baseUrl + kycId.String())
	//
	//		if err != nil {
	//			Expect(err).NotTo(HaveOccurred())
	//		}
	//
	//		By("status code should be 200")
	//		Expect(resp.StatusCode()).To(Equal(http.StatusOK))
	//	})
	//})
})
