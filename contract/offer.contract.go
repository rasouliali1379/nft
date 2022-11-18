package contract

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"nft/infra/persist/type"
	"nft/internal/offer/model"
)

type IOfferController interface {
	MakeOffer(c *fiber.Ctx) error
	CancelOffer(c *fiber.Ctx) error
	AcceptOffer(c *fiber.Ctx) error
	GetAllOffers(c *fiber.Ctx) error
}

type IOfferService interface {
	MakeOfferToSale(c context.Context, m model.Offer) error
	CancelOffer(c context.Context, m model.Offer) error
	AcceptOffer(c context.Context, m model.Offer) error
	GetAllOffers(c context.Context, m model.Offer) ([]model.Offer, error)
}

type IOfferRepository interface {
	Add(c context.Context, m model.Offer) (model.Offer, error)
	Get(c context.Context, conditions persist.D) (model.Offer, error)
	GetAll(c context.Context, conditions persist.D) ([]model.Offer, error)
	Delete(c context.Context, m model.Offer) error
	Update(c context.Context, m model.Offer) (model.Offer, error)
}
