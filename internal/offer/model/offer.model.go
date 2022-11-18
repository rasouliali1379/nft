package model

import (
	"github.com/google/uuid"
	user "nft/internal/user/model"
)

type Offer struct {
	ID       *uuid.UUID
	User     user.User
	SaleId   uuid.UUID
	Price    float64
	Accepted bool
}
