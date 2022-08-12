package kyc

import (
	file "nft/src/file/model"

	"github.com/google/uuid"
)

type KYC struct {
	ID              uuid.UUID
	ApprovedBy      *uuid.UUID
	RejectedBy      *uuid.UUID
	RejectionReason string
	IdCardImage     file.Image
	PortraitImage   file.Image
	UserId          uuid.UUID
}
