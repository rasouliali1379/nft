package contract

import (
	"context"
	"github.com/gofiber/fiber/v2"
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
}

type ISaleRepository interface {
	Create(c context.Context, m model.Sale) (model.Sale, error)
}
