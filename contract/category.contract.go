package contract

import (
	"context"
	"nft/infra/persist/type"
	model "nft/internal/category/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ICategoryController interface {
	AddCategory(c *fiber.Ctx) error
	GetCategory(c *fiber.Ctx) error
	GetAllCategories(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type ICategoryService interface {
	GetCategory(c context.Context, id uuid.UUID) (model.Category, error)
	GetAllCategories(c context.Context) ([]model.Category, error)
	GetSubCategories(c context.Context, id uuid.UUID) ([]model.Category, error)
	AddCategory(c context.Context, category model.Category) (model.Category, error)
	UpdateCategory(c context.Context, category model.Category) (model.Category, error)
	DeleteCategory(c context.Context, catId uuid.UUID) error
}

type ICategoryRepository interface {
	Exists(c context.Context, conditions persist.D) error
	Add(c context.Context, category model.Category) (model.Category, error)
	Update(c context.Context, userModel model.Category) (model.Category, error)
	Delete(c context.Context, userId uuid.UUID) error
	Get(c context.Context, conditions persist.D) (model.Category, error)
	GetAll(c context.Context, conditions persist.D) ([]model.Category, error)
}
