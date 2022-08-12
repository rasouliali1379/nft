package kyc

import "github.com/google/uuid"

type KYC struct {
	ID              uuid.UUID `json:"id"`
	IdCardImageUrl  string    `json:"id_card_image_url"`
	PortraitImage   string    `json:"portrait_image"`
	Status          KYCStatus `json:"status"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
}

type KYCStatus string

const (
	KYCStatusApproved  KYCStatus = "approved"
	KYCStatusRejected  KYCStatus = "rejected"
	KYCStatusUndefined KYCStatus = "undefined"
)

type KYCList struct {
	KYCList []KYC `json:"kyc_list"`
}

type RejectAppeal struct {
	Message string `json:"message"`
}
