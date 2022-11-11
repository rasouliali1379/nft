package contract

import (
	"context"
	"nft/infra/persist/model"
	model "nft/internal/card/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ICardController interface {
	GetCard(c *fiber.Ctx) error
	GetAllCards(c *fiber.Ctx) error
	AddCard(c *fiber.Ctx) error
	RemoveCard(c *fiber.Ctx) error
	ApproveCard(c *fiber.Ctx) error
}

type ICardService interface {
	GetAllCards(c context.Context, userId uuid.UUID) ([]model.Card, error)
	GetCard(c context.Context, id uuid.UUID, userId uuid.UUID) (model.Card, error)
	AddCard(c context.Context, cardModel model.Card) (model.Card, error)
	ApproveCard(c context.Context, id uuid.UUID, userId uuid.UUID) error
	DeleteCard(c context.Context, id uuid.UUID, userId uuid.UUID) error
}

type ICardRepository interface {
	Exists(c context.Context, conditions persist.Conds) (bool, error)
	Add(c context.Context, cardModel model.Card) (model.Card, error)
	Update(c context.Context, cardModel model.Card) (model.Card, error)
	Delete(c context.Context, cardId uuid.UUID) error
	Get(c context.Context, conditions persist.Conds) (model.Card, error)
	GetAll(c context.Context, conditions persist.Conds) ([]model.Card, error)
}
