package category

import (
	"github.com/google/uuid"
	usermodel "nft/internal/user/model"
	"time"
)

type Category struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	Name          string
	SubCategories []Category
	ParentId      *uuid.UUID
	CreatedBy     usermodel.User
}
