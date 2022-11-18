package contract

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"nft/infra/persist/type"
	model "nft/internal/collection/model"
)

type ICollectionController interface {
	Add(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type ICollectionService interface {
	GetCollection(c context.Context, m model.Collection) (model.Collection, error)
	GetAllCollections(c context.Context, query model.QueryCollection) ([]model.Collection, error)
	AddCollection(c context.Context, m model.Collection) (model.Collection, error)
	DeleteCollection(c context.Context, m model.Collection) error
	GetOwnedCollection(c context.Context, m model.Collection) (model.Collection, error)
}

type ICollectionRepository interface {
	Add(c context.Context, category model.Collection) (model.Collection, error)
	Update(c context.Context, userModel model.Collection) (model.Collection, error)
	Delete(c context.Context, m model.Collection) error
	Get(c context.Context, conditions persist.D) (model.Collection, error)
	GetAll(c context.Context, conditions persist.D) ([]model.Collection, error)
	HardDelete(c context.Context, id uuid.UUID) error
}
