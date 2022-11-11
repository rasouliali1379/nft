package offer

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"nft/contract"
)

type OfferController struct {
	offerService contract.IOfferService
}

type OfferControllerParams struct {
	fx.In
	OfferService contract.IOfferService
}

func NewOfferController(params OfferControllerParams) contract.IOfferController {
	return &OfferController{
		offerService: params.OfferService,
	}
}

func (o OfferController) MakeOffer(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (o OfferController) CancelOffer(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (o OfferController) AcceptOffer(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (o OfferController) RejectOffer(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (o OfferController) GetAllOffers(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
