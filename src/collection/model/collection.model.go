package collection

import (
	"github.com/google/uuid"
	catmodel "nft/src/category/model"
	file "nft/src/file/model"
	usermodel "nft/src/user/model"
)

type Collection struct {
	ID          *uuid.UUID
	Title       string
	Description string
	Categories  []catmodel.Category
	User        usermodel.User
	HeaderImage *file.Image
	Status      CollectionStatus
}

type CollectionStatus string

const (
	CollectionStatusDraft CollectionStatus = "Draft"
	CollectionStatusSaved                  = "Saved"
)
