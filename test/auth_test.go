package test

import (
	"encoding/json"
	"fmt"
	"maskan/config"
	authdto "maskan/src/auth/dto"
	"net/http"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	client := resty.New()
	baseUrl := fmt.Sprintf("http://%s:%s/v1/auth/", config.C().App.Http.Host, config.C().App.Http.Port)
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

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Login", func() {
		It("should login the user successfully", func() {
			resp, err := client.R().
				SetBody(loginDto).
				Post(baseUrl + "login")

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
