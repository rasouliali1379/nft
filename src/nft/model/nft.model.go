package nft

import (
	"github.com/google/uuid"
	catmodel "nft/src/category/model"
	file "nft/src/file/model"
	user "nft/src/user/model"
)

type Nft struct {
	ID          *uuid.UUID
	Title       string
	Description string
	//CollectionId *uuid.UUID
	Categories      []catmodel.Category
	User            user.User
	Status          NftStatus
	NftImage        *file.Image
	RejectedBy      *user.User
	RejectionReason string
	ApprovedBy      *user.User
}

type NftStatus string

const (
	NftStatusDraft     NftStatus = "Draft"
	NftStatusPending             = "Pending"
	NftStatusRejected            = "Rejected"
	NftStatusApproved            = "Approved"
	NftStatusProcessed           = "Processed"
)
