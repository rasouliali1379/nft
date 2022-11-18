package offer

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log"
	"net/http"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	"nft/internal/offer/dto"
	"nft/internal/offer/model"
	usermodel "nft/internal/user/model"
	"nft/pkg/filper"
	"nft/pkg/validator"
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

// MakeOffer godoc
// @Summary  make offer on auction or p2p sale
// @Tags     offer
// @Accept   json
// @Produce  json
// @Router   /v1/offer/ [post]
// @Param    message  body      dto.MakeOfferRequest  true  "price should be higher than sale min price"
// @Success  200      {string}  string                "offer made successfully"
func (o OfferController) MakeOffer(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "OfferController[MakeOffer]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.MakeOfferRequest
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errRes := validator.Validate(request)
	if len(errRes.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	offer := model.Offer{
		User:   usermodel.User{ID: userId},
		SaleId: request.SaleId,
		Price:  request.Price,
	}
	if err := o.offerService.MakeOfferToSale(ctx, offer); err != nil {
		log.Println(err)
		if errors.Is(err, apperrors.ErrSaleNotFound) {
			return filper.GetNotFoundError(c, "sale not found")
		} else if errors.Is(err, apperrors.ErrOfferYourSale) {
			return filper.GetNotFoundError(c, "you can't make offer on your own sale")
		} else if errors.Is(err, apperrors.ErrOfferLowerMinPrice) {
			return filper.GetNotFoundError(c, "your offer should be higher than sale min price")
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "offer made successfully")
}

// CancelOffer godoc
// @Summary  remove offer from auction or p2p sale
// @Tags     offer
// @Accept   json
// @Produce  json
// @Param    id  path  string  true  "offer id that you want to be canceled"
// @Router   /v1/offer/{id} [delete]
// @Success  200  {string}  string  "offer canceled successfully"
func (o OfferController) CancelOffer(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "OfferController[CancelOffer]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	offerId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid offer id")
	}

	if err := o.offerService.CancelOffer(ctx, model.Offer{ID: &offerId, User: usermodel.User{ID: userId}}); err != nil {
		return err
	}

	return filper.GetSuccessResponse(c, "offer canceled successfully")
}

// AcceptOffer godoc
// @Summary  accept offer for p2p sales only
// @Tags     offer
// @Accept   json
// @Produce  json
// @Param    id  path  string  true  "offer id that you want to accept"
// @Router   /v1/offer/{id}/accept [post]
// @Success  200  {string}  string  "offer accepted successfully"
func (o OfferController) AcceptOffer(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "OfferController[AcceptOffer]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	offerId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid offer id")
	}

	if err := o.offerService.AcceptOffer(ctx, model.Offer{ID: &offerId, User: usermodel.User{ID: userId}}); err != nil {
		if errors.Is(err, apperrors.ErrOfferNotFound) {
			return filper.GetNotFoundError(c, "offer not found")
		} else if errors.Is(err, apperrors.ErrSaleNotFound) {
			return filper.GetNotFoundError(c, "sale not found")
		}
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "offer accepted successfully")
}

// GetAllOffers godoc
// @Summary  get sale offer list
// @Tags     offer
// @Accept   json
// @Produce  json
// @Router   /v1/offer [get]
// @Success  200      {object}  dto.OfferList
// @Param    sale_id  query     string  true  "sale id that you want its offers"  Format(uuid)
func (o OfferController) GetAllOffers(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "OfferController[GetAllOffers]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	if len(c.Query("sale_id")) == 0 {
		return filper.GetBadRequestError(c, "sale id is mandatory")
	}
	saleId, err := uuid.Parse(c.Query("sale_id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid sale id")
	}
	offerList, err := o.offerService.GetAllOffers(ctx, model.Offer{SaleId: saleId, User: usermodel.User{ID: userId}})
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.Status(http.StatusOK).JSON(createOfferListDtoFromModel(offerList))
}
