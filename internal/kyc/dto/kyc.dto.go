package kyc

import "github.com/google/uuid"

type Kyc struct {
	ID              uuid.UUID `json:"id"`
	IdCardImageUrl  string    `json:"id_card_image_url"`
	PortraitImage   string    `json:"portrait_image"`
	Status          KycStatus `json:"status"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
}

type KycStatus string

const (
	KYCStatusApproved  KycStatus = "approved"
	KYCStatusRejected  KycStatus = "rejected"
	KYCStatusUndefined KycStatus = "undefined"
)

type KycList struct {
	KYCList []Kyc `json:"kyc_list"`
}

type RejectAppeal struct {
	Message string `json:"message"`
}
