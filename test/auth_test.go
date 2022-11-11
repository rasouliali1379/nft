package test

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"nft/config"
	authdto "nft/internal/auth/dto"
	jwt "nft/internal/jwt/model"
)

var _ = Describe("Auth", Ordered, func() {

	var baseUrl string

	client := resty.New()
	signUpDto := authdto.SignUpRequest{
		FirstName:      "Ali",
		LastName:       "Rasouli",
		NationalId:     "0123456789",
		Email:          "ali3@gmail.com",
		PhoneNumber:    "09368045731",
		LandLineNumber: "02133073333",
		Password:       "ali1379",
		Province:       "tehran",
		City:           "terhan",
		Address:        "mahallati",
	}

	loginDto := authdto.LoginRequest{
		Email:    "ali3@gmail.com",
		Password: "ali1379",
	}

	var otpToken string
	var jwtToken jwt.Jwt

	BeforeAll(func() {
		baseUrl = fmt.Sprintf("http://%s:%s/v1/auth/", config.C().App.Http.Host, config.C().App.Http.Port)
	})

	Describe("SignUp", func() {
		It("should sign up new user successfully", func() {
			resp, err := client.R().
				SetBody(signUpDto).
				Post(baseUrl + "signup")

			Expect(err).NotTo(HaveOccurred())

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

			var signUpResponse authdto.OtpToken
			err = json.Unmarshal(resp.Body(), &signUpResponse)
			Expect(err).NotTo(HaveOccurred())

			otpToken = signUpResponse.Token

			By("otp token shouldn't be empty")
			Expect(resp.StatusCode()).NotTo(Equal(""))
		})
	})

	Describe("Verify Email", func() {
		It("should verify email successfully", func() {

			resp, err := client.R().
				SetBody(authdto.VerifyEmailRequest{
					Token: otpToken,
					Code:  "111111",
				}).
				Post(baseUrl + "verify-email")
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Resend Verification Email", func() {
		It("should resend verification email successfully", func() {

			resp, err := client.R().
				SetBody(authdto.ResendEmailRequest{Token: otpToken}).
				Post(baseUrl + "resend-email")
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Login", func() {
		It("should login the user successfully", func() {

			resp, err := client.R().
				SetBody(loginDto).
				Post(baseUrl + "login")
			Expect(err).NotTo(HaveOccurred())

			if err = json.Unmarshal(resp.Body(), &jwtToken); err != nil {
				Fail("unable to unmarshal jwt object")
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Refresh", func() {
		It("should refresh user token successfully", func() {

			resp, err := client.R().
				SetBody(authdto.RefreshRequest{RefreshToken: jwtToken.RefreshToken}).
				Post(baseUrl + "refresh")
			Expect(err).NotTo(HaveOccurred())

			if err = json.Unmarshal(resp.Body(), &jwtToken); err != nil {
				Fail("unable to unmarshal jwt object")
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Logout", func() {
		It("should logout the user successfully", func() {

			resp, err := client.R().
				SetBody(authdto.RefreshRequest{RefreshToken: jwtToken.RefreshToken}).
				Post(baseUrl + "logout")
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
