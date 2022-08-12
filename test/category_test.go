package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nft/config"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	ng "github.com/goombaio/namegenerator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	authdto "nft/src/auth/dto"
	catdto "nft/src/category/dto"
	jwt "nft/src/jwt/model"
)

var token string

var _ = Describe("Category Management", Ordered, func() {

	BeforeAll(func() {
		client := resty.New()
		baseUrl := fmt.Sprintf("http://%s:%s/v1/auth/", config.C().App.Http.Host, config.C().App.Http.Port)

		signUpDto := authdto.SignUpRequest{
			FirstName:      "Ali",
			LastName:       "Rasouli",
			NationalId:     "0123456782",
			Email:          "testcat@gmail.com",
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
			Email:    "testcat@gmail.com",
			Password: "ali1379",
		}

		resp, err = client.R().
			SetBody(loginDto).
			Post(baseUrl + "login")

		if err != nil {
			AbortSuite(fmt.Sprintf("failed to unmarshal login response: %s", err.Error()))
		}

		var jwtToken jwt.Jwt
		err = json.Unmarshal(resp.Body(), &jwtToken)
		if err != nil {
			AbortSuite(fmt.Sprintf("failed unmarshal jwt struct: %s", err.Error()))
		}
		token = jwtToken.AccessToken
	})

	cat := catdto.AddCategoryRequest{
		Name: "First",
	}

	var catId uuid.UUID
	var cat2 catdto.CategoryDto

	client := resty.New()
	baseUrl := fmt.Sprintf("http://%s:%s/v1/category/", config.C().App.Http.Host, config.C().App.Http.Port)

	Describe("add new category", func() {
		It("should add new category successfully", func() {

			resp, err := client.R().
				SetBody(cat).
				SetAuthToken(token).
				Post(baseUrl)

			Expect(err).NotTo(HaveOccurred())

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

			var cat1 catdto.CategoryDto
			err = json.Unmarshal(resp.Body(), &cat1)
			Expect(err).NotTo(HaveOccurred())

			catId = cat1.ID
		})

		It("should add new subcategory successfully", func() {
			resp, err := client.R().
				SetBody(catdto.AddCategoryRequest{
					Name:     "Fisrt",
					ParentId: catId,
				}).
				SetAuthToken(token).
				Post(fmt.Sprintf("%s/", baseUrl))
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

			err = json.Unmarshal(resp.Body(), &cat2)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Get categories list", func() {
		It("should get users list successfully", func() {
			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl)

			if err != nil {
				Fail(fmt.Sprintf("failed to make request to get categories list: %s", err.Error()))
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var catList catdto.CategoriesListDto
			if err = json.Unmarshal(resp.Body(), &catList); err != nil {
				Fail("unable to unmarshal user list")
			}

			By("categories list should have one item")
			Expect(len(catList.Categories)).To(Equal(1))
		})
	})

	Describe("Get category", func() {
		It("should get single category successfully", func() {

			if catId == uuid.Nil {
				Fail("category id is empty")
			}

			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl + catId.String())

			if err != nil {
				Fail(fmt.Sprintf("failed to make the request to get category: %s", err.Error()))
			}

			var cat1 catdto.CategoryDto
			if err = json.Unmarshal(resp.Body(), &cat1); err != nil {
				Fail(fmt.Sprintf("unable to unmarshal user details: %s", err.Error()))
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			By("should have one subcategory")
			Expect(len(cat1.SubCategories)).To(Equal(1))
		})
	})

	Describe("Update category", func() {
		It("should update category successfully", func() {

			generatedName := ng.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()

			resp, err := client.R().
				SetBody(catdto.AddCategoryRequest{
					Name: generatedName,
				}).
				SetAuthToken(token).
				Patch(baseUrl + catId.String())

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			var updatedCat catdto.CategoryDto

			if err = json.Unmarshal(resp.Body(), &updatedCat); err != nil {
				Fail("unable to unmarshal user details")
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			By(fmt.Sprintf("returned category name should be %s", generatedName))
			Expect(updatedCat.Name).To(Equal(generatedName))
		})
	})

	Describe("Delete category", func() {
		It("should delete user successfully", func() {

			resp, err := client.R().
				SetAuthToken(token).
				Delete(baseUrl + catId.String())

			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
