package card

import (
	"errors"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	dto "nft/internal/card/dto"
	model "nft/internal/card/model"
	"nft/pkg/filper"
	"nft/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type CardController struct {
	CardService contract.ICardService
}

type CardControllerParams struct {
	fx.In
	CardService contract.ICardService
}

func NewCardController(params CardControllerParams) contract.ICardController {
	return &CardController{
		CardService: params.CardService,
	}
}

// GetCard godoc
// @Summary  get single card
// @Tags     card
// @Accept   json
// @Produce  json
// @Param    id   path      int  true  "card id that will be retrieved"
// @Success  200  {object}  dto.Card
// @Router   /v1/card/{id} [get]
func (cc CardController) GetCard(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CardController[GetCard]")
	defer span.Finish()

	userId := c.Locals("user_id").(uuid.UUID)

	cardId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid card id")
	}

	catModel, err := cc.CardService.GetCard(ctx, cardId, userId)
	if err != nil {
		if errors.Is(err, apperrors.ErrCardDoesntBelongToUser) {
			return filper.GetNotFoundError(c, "card doesn't exists")
		}
		return filper.GetInternalError(c, "")
	}

	return c.JSON(mapCardModelToDTo(catModel))
}

// GetAllCards godoc
// @Summary  get cards list
// @Tags     card
// @Accept   json
// @Produce  json
// @Success  200  {object}  dto.CardList
// @Router   /v1/card [get]
func (cc CardController) GetAllCards(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CardController[GetAllCards]")
	defer span.Finish()

	userId := c.Locals("user_id").(uuid.UUID)

	cardsList, err := cc.CardService.GetAllCards(ctx, userId)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.JSON(mapCardListModelToDto(cardsList))
}

// AddCard godoc
// @Summary  add new card
// @Tags     card
// @Accept   json
// @Produce  json
// @Param    message  body      dto.AddCardRequest  true  "Add card request body."
// @Success  200      {object}  dto.Card
// @Router   /v1/card [post]
func (cc CardController) AddCard(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CardController[AddCard]")
	defer span.Finish()

	userId := c.Locals("user_id").(uuid.UUID)

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.AddCardRequest
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errRes := validator.Validate(request)
	if len(errRes.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)

	}

	cardModel, err := cc.CardService.AddCard(ctx, model.Card{
		CardNumber: request.CardNumber,
		IBAN:       request.IBAN,
		UserId:     userId,
	})

	if err != nil {
		if errors.Is(err, apperrors.ErrCardAlreadyExistsForUser) {
			return filper.GetBadRequestError(c, "card already exists for user")
		}
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapCardModelToDTo(cardModel))
}

// RemoveCard godoc
// @Summary  remove existing card
// @Tags     card
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "card id that will be removed"
// @Success  200  {string}  string  "card removed successfully"
// @Router   /v1/card/{id} [delete]
func (cc CardController) RemoveCard(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CardController[RemoveCard]")
	defer span.Finish()

	userId := c.Locals("user_id").(uuid.UUID)

	cardId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid card id")
	}

	err = cc.CardService.DeleteCard(ctx, cardId, userId)
	if err != nil {
		if errors.Is(err, apperrors.ErrCardDoesntBelongToUser) {
			return filper.GetNotFoundError(c, "card wasn't found")
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "card removed successfully")
}

// ApproveCard godoc
// @Summary  approve card
// @Tags     card
// @Accept   json
// @Produce  json
// @Param    id   path      int     true  "card id that will be approved"
// @Success  200  {string}  string  "card approved successfully"
// @Router   /v1/card/{id}/approve [get]
func (cc CardController) ApproveCard(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "CardController[ApproveCard]")
	defer span.Finish()

	userId := c.Locals("user_id").(uuid.UUID)
	cardId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid card id")
	}

	if err := cc.CardService.ApproveCard(ctx, cardId, userId); err != nil {
		if errors.Is(err, apperrors.ErrCardDoesntBelongToUser) {
			return filper.GetNotFoundError(c, "card wasn't found")
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "card approved successfully")
}
