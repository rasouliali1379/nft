package contract

import "github.com/gofiber/fiber/v2"

type IOfferController interface {
	MakeOffer(c *fiber.Ctx) error
	CancelOffer(c *fiber.Ctx) error
	AcceptOffer(c *fiber.Ctx) error
	RejectOffer(c *fiber.Ctx) error
	GetAllOffers(c *fiber.Ctx) error
}

type IOfferService interface {
}

type IOfferRepository interface {
}
