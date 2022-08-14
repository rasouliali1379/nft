package test

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"nft/config"
	dto "nft/src/kyc/dto"
	"os"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kyc Management", Ordered, func() {
	var kycImage1, kycImage2 *os.File
	var baseUrl string
	var kycId uuid.UUID

	client := resty.New()

	BeforeAll(func() {
		baseUrl = fmt.Sprintf("http://%s:%s/v1/kyc/", config.C().App.Http.Host, config.C().App.Http.Port)
		path, err := os.Getwd()
		if err != nil {
			Fail(fmt.Sprintf("failed get current directory path: %s", err.Error()), 7)
		}

		kycImage1, err = os.Open(path + "/assets/kyc-test-image-1.jpeg")
		if err != nil {
			Fail(fmt.Sprintf("failed to open kyc image file 1: %s", err.Error()), 7)
		}

		kycImage2, err = os.Open(path + "/assets/kyc-test-image-2.jpg")
		if err != nil {
			Fail(fmt.Sprintf("failed to open kyc image file 2: %s", err.Error()), 7)
		}
	})

	Describe("appeal for kyc", func() {
		It("should request for kyc successfully", func() {
			resp, err := client.R().
				SetMultipartFields(
					&resty.MultipartField{
						Param:    "id_card",
						Reader:   kycImage1,
						FileName: kycImage1.Name(),
					},
					&resty.MultipartField{
						Param:    "portrait",
						Reader:   kycImage2,
						FileName: kycImage2.Name(),
					}).
				SetAuthToken(token).
				Post(baseUrl)

			if err != nil {
				Fail(fmt.Sprintf("failed appeal for kyc: %s", err.Error()), 6)
			}

			By("status code should be 201")
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
		})
	})

	Describe("Get kyc list", func() {
		It("should get kyc list successfully", func() {
			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl)

			if err != nil {
				Fail(fmt.Sprintf("failed to make request to get kycs list: %s", err.Error()), 3)
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var kycList dto.KYCList
			if err = json.Unmarshal(resp.Body(), &kycList); err != nil {
				Fail("unable to unmarshal kyc list", 3)
			}

			By("kyc list should have one item")
			Expect(len(kycList.KYCList)).To(Equal(1))

			kycId = kycList.KYCList[0].ID
		})
	})

	Describe("Get kyc", func() {
		It("should get single kyc successfully", func() {

			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl + kycId.String())

			if err != nil {
				Fail(fmt.Sprintf("failed to make the request to get kyc: %s", err.Error()))
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var kyc dto.KYC
			if err = json.Unmarshal(resp.Body(), &kyc); err != nil {
				Fail("unable to unmarshal kyc object")
			}

			By("status field should be equal to undefined")
			Expect(kyc.Status).To(Equal(dto.KYCStatusUndefined))
		})
	})

	Describe("Approve kyc", func() {
		It("should approve kyc successfully", func() {
			resp, err := client.R().
				SetAuthToken(token).
				Post(baseUrl + kycId.String() + "/approve")

			if err != nil {
				Fail(fmt.Sprintf("failed to make the request to approve kyc: %s", err.Error()), 1)
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})

		It("should get single kyc and check its status to be approved", func() {

			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl + kycId.String())

			if err != nil {
				Fail(fmt.Sprintf("failed to make the request to get kyc: %s", err.Error()))
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var kyc dto.KYC
			if err = json.Unmarshal(resp.Body(), &kyc); err != nil {
				Fail("unable to unmarshal kyc object")
			}

			By("status field should be equal to approved")
			Expect(kyc.Status).To(Equal(dto.KYCStatusApproved))
		})
	})

	Describe("Reject kyc", func() {
		It("should reject kyc successfully", func() {
			resp, err := client.R().
				SetAuthToken(token).
				SetBody(dto.RejectAppeal{Message: "some reason"}).
				Post(baseUrl + kycId.String() + "/reject")

			if err != nil {
				Fail(fmt.Sprintf("failed to make the request to reject kyc: %s", err.Error()), 1)
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
		})

		It("should get single kyc and check its status to be rejected", func() {
			resp, err := client.R().
				SetAuthToken(token).
				Get(baseUrl + kycId.String())

			if err != nil {
				Fail(fmt.Sprintf("failed to make the request to get kyc: %s", err.Error()))
			}

			By("status code should be 200")
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var kyc dto.KYC
			if err = json.Unmarshal(resp.Body(), &kyc); err != nil {
				Fail("unable to unmarshal kyc object")
			}

			By("status field should be equal to rejected")
			Expect(kyc.Status).To(Equal(dto.KYCStatusRejected))
		})
	})
})
