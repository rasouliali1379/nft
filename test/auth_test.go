package test

import (
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

	Describe("SignUp", func() {
		It("should sign up new user successfully", func() {
			resp, err := client.R().
				SetBody(signUpDto).
				Post(baseUrl + "signup")

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
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
