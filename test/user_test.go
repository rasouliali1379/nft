package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nft/config"
	"time"

	"github.com/go-resty/resty/v2"
	ng "github.com/goombaio/namegenerator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	authdto "nft/src/auth/dto"
	userdto "nft/src/user/dto"
)

var _ = Describe("User Management", Ordered, func() {
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
	var baseUrl string
	client := resty.New()

	BeforeAll(func() {
		baseUrl = fmt.Sprintf("http://%s:%s/v1/user/", config.C().App.Http.Host, config.C().App.Http.Port)
	})

	Describe("add new user", func() {
		It("should add new user successfully", func() {
			resp, err := client.R().
				SetBody(user).
				Post(baseUrl)
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
		})
	})

	Describe("Get users list", func() {
		It("should get users list successfully", func() {

			resp, err := client.R().
				Get(baseUrl)
			if err != nil {
				Fail(fmt.Sprintf("unable to make request to get user list: %s", err.Error()), 3)
			}

			err = json.Unmarshal(resp.Body(), &userList)
			if err != nil {
				Fail(fmt.Sprintf("unable to unmarshal user list: %s", err.Error()), 3)
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})

	Describe("Get user", func() {
		It("should get single user successfully", func() {

			resp, err := client.R().
				Get(baseUrl + userList.Users[0].ID)
			Expect(err).NotTo(HaveOccurred())

			var user userdto.UserDto
			err = json.Unmarshal(resp.Body(), &user)
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			By(fmt.Sprintf("returned user should be %s", userList.Users[0]))
			Expect(user).To(Equal(userList.Users[0]))
		})
	})

	Describe("Update user", func() {
		It("should update user successfully", func() {

			userDetails := userList.Users[0]
			generatedName := ng.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
			userDetails.FirstName = generatedName

			resp, err := client.R().
				SetBody(userDetails).
				Patch(baseUrl + userDetails.ID)
			Expect(err).NotTo(HaveOccurred())

			var updatedUser userdto.UserDto
			err = json.Unmarshal(resp.Body(), &updatedUser)
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			By(fmt.Sprintf("returned user first name should be %s", generatedName))
			Expect(updatedUser.FirstName).To(Equal(generatedName))
		})
	})

	Describe("Delete user", func() {
		It("should delete user successfully", func() {

			resp, err := client.R().
				Delete(baseUrl + userList.Users[0].ID)
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
