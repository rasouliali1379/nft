package test

import (
	"encoding/json"
	"fmt"
	"maskan/config"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	ng "github.com/goombaio/namegenerator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	authdto "maskan/src/auth/dto"
	userdto "maskan/src/user/dto"
)

var _ = Describe("User Management", func() {
	user := authdto.SignUpRequest{
		FirstName:      "Mohammad",
		LastName:       "Javadi",
		NationalId:     "2231332323",
		Email:          "samrd200046@gmail.com",
		PhoneNumber:    "09222222222",
		LandLineNumber: "33232323",
		Province:       "Tehran",
		City:           "Tehran",
		Address:        "mahallati",
		Password:       "mohammad1372",
	}

	var userList userdto.UserListDto

	client := resty.New()
	baseUrl := fmt.Sprintf("http://%s:%s/v1/user/", config.C().App.Http.Host, config.C().App.Http.Port)

	Describe("add new user", func() {
		It("should add new user successfully", func() {
			resp, err := client.R().
				SetBody(user).
				Post(fmt.Sprintf("%s/", baseUrl))

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
		})
	})

	Describe("Get users list", func() {
		It("should get users list successfully", func() {
			resp, err := client.R().
				Get(baseUrl)

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			if err = json.Unmarshal(resp.Body(), &userList); err != nil {
				Fail("unable to unmarshal user list")
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Get user", func() {
		It("should get users list successfully", func() {
			resp, err := client.R().
				Get(baseUrl + userList.Users[0].ID)

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			var user userdto.UserDto

			if err = json.Unmarshal(resp.Body(), &user); err != nil {
				Fail("unable to unmarshal user details")
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			By(fmt.Sprintf("returned user should be %s", userList.Users[0]))
			Expect(user).To(Equal(userList.Users[0]))
		})
	})

	Describe("Update user", func() {
		It("should update user successfully", func() {

			if userList.Users == nil {
				Fail("user list is empty")
			}

			userDetails := userList.Users[0]
			generatedName := ng.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
			userDetails.FirstName = generatedName

			resp, err := client.R().
				SetBody(userDetails).
				Patch(baseUrl + userDetails.ID)

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			var updatedUser userdto.UserDto

			if err = json.Unmarshal(resp.Body(), &updatedUser); err != nil {
				Fail("unable to unmarshal user details")
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			By(fmt.Sprintf("returned user first name should be %s", generatedName))
			Expect(updatedUser.FirstName).To(Equal(generatedName))
		})
	})

	Describe("Delete user", func() {
		It("should delete user successfully", func() {

			if userList.Users == nil {
				Fail("user list is empty")
			}

			userDetails := userList.Users[0]

			resp, err := client.R().
				SetBody(userDetails).
				Delete(baseUrl + userDetails.ID)

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
