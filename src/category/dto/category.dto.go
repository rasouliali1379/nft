package category

import (
	"github.com/google/uuid"
	userdto "nft/src/user/dto"
	"time"
)

type CategoryDto struct {
	ID        uuid.UUID  `json:"id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Name          string          `json:"name,omitempty"`
	SubCategories []CategoryDto   `json:"sub_categories,omitempty"`
	CreatedBy     userdto.UserDto `json:"created_by,omitempty"`
}

type CategoriesListDto struct {
	Categories []CategoryDto `json:"categories"`
}

type AddCategoryRequest struct {
	Name     string    `json:"name" validate:"required"`
	ParentId uuid.UUID `json:"parent_id"`
}

type AddCategoryResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
}
