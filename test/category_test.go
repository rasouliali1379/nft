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

	catdto "nft/src/category/dto"
)

var _ = Describe("Category Management", Ordered, func() {
	cat := catdto.AddCategoryRequest{
		Name: "First",
	}

	var catId uuid.UUID
	var cat2 catdto.CategoryDto
	var baseUrl string
	client := resty.New()

	BeforeAll(func() {
		baseUrl = fmt.Sprintf("http://%s:%s/v1/category/", config.C().App.Http.Host, config.C().App.Http.Port)
	})

	Describe("add new category", func() {
		It("should add new category successfully", func() {

			resp, err := client.R().
				SetBody(cat).
				SetAuthToken(token).
				Post(baseUrl)
			if err != nil {
				Fail(fmt.Sprintf("unable to make request to create category: %s", err.Error()), 5)
			}

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

			var cat1 catdto.CategoryDto
			err = json.Unmarshal(resp.Body(), &cat1)
			if err != nil {
				Fail(fmt.Sprintf("unable to unmarshal category: %s", err.Error()), 5)
			}

			catId = cat1.ID
		})

		It("should add new subcategory successfully", func() {

			resp, err := client.R().
				SetBody(catdto.AddCategoryRequest{
					Name:     "First",
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
		It("should get category list successfully", func() {

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
				Fail("unable to unmarshal category list")
			}

			By("categories list should have one item")
			Expect(len(catList.Categories)).To(Equal(1))
		})
	})

	Describe("Get category", func() {
		It("should get single category successfully", func() {

			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl + catId.String())
			Expect(err).NotTo(HaveOccurred())

			var cat1 catdto.CategoryDto
			err = json.Unmarshal(resp.Body(), &cat1)
			Expect(err).NotTo(HaveOccurred())

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
			Expect(err).NotTo(HaveOccurred())

			var updatedCat catdto.CategoryDto
			err = json.Unmarshal(resp.Body(), &updatedCat)
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			By(fmt.Sprintf("returned category name should be %s", generatedName))
			Expect(updatedCat.Name).To(Equal(generatedName))
		})
	})

	Describe("Delete category", func() {
		It("should delete category successfully", func() {

			resp, err := client.R().
				SetAuthToken(token).
				Delete(baseUrl + catId.String())
			Expect(err).NotTo(HaveOccurred())

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})
	})
})
