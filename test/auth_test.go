package test

import (
	"fmt"
	authdto "maskan/src/auth/dto"
	"net/http"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	client := resty.New()
	baseUrl := "http://localhost:1212/v1/auth"
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
		It("should register the user successfully", func() {
			resp, _ := client.R().
				SetBody(signUpDto).
				Post(fmt.Sprintf("%s/signup", baseUrl))

			// if err != nil {
			// 	log.Println(err)
			// }
			// log.Println(resp)

			By("The status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
		})
	})

	Describe("Login", func() {
		It("should register the user successfully", func() {
			resp, _ := client.R().
				SetBody(loginDto).
				Post(fmt.Sprintf("%s/login", baseUrl))

			// if err != nil {
			// 	log.Println(err)
			// }
			// log.Println(resp)

			By("The status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
