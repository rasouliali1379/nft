package contract

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"nft/infra/persist/type"
	"nft/internal/sale/model"
)

type ISaleController interface {
	SellNft(c *fiber.Ctx) error
	SellCollection(c *fiber.Ctx) error
	CancelSale(c *fiber.Ctx) error
	GetAllSales(c *fiber.Ctx) error
	GetSale(c *fiber.Ctx) error
}

type ISaleService interface {
	CreateNftSale(c context.Context, m model.Sale) (model.Sale, error)
	CreateCollectionSale(c context.Context, m model.Sale) (model.Sale, error)
	CancelSale(c context.Context, m model.Sale) error
	GetSalesList(c context.Context, userId uuid.UUID) ([]model.Sale, error)
	GetSale(c context.Context, m model.Sale) (model.Sale, error)
}

type ISaleRepository interface {
	Create(c context.Context, m model.Sale) (model.Sale, error)
	Get(c context.Context, conditions persist.D) (model.Sale, error)
	GetAll(c context.Context, conditions persist.D) ([]model.Sale, error)
	Cancel(c context.Context, m model.Sale) error
}
