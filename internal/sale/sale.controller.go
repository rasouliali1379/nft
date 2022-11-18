package sale

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log"
	"nft/contract"
	apperrors "nft/error"
	"nft/infra/jtrace"
	dto "nft/internal/sale/dto"
	"nft/internal/sale/model"
	usermodel "nft/internal/user/model"
	"nft/pkg/filper"
	"nft/pkg/validator"
)

type SaleController struct {
	saleService contract.ISaleService
}

type SaleControllerParams struct {
	fx.In
	SaleService contract.ISaleService
}

func NewSaleController(params SaleControllerParams) contract.ISaleController {
	return &SaleController{
		saleService: params.SaleService,
	}
}

// SellNft godoc
// @Summary  sell or hold auction for nft assets
// @Tags     sale
// @Accept   json
// @Produce  json
// @Router   /v1/sale/sell-nft [post]
// @Param    message  body  dto.SaleRequest  true  "nft sale request body"
// @Success  200      body  model.Sale
func (s SaleController) SellNft(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "SaleController[SellNft]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.SaleRequest
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errRes := validator.Validate(request)
	if len(errRes.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	sale, err := s.saleService.CreateNftSale(ctx, mapCreateSaleDtoToModel(request, userId))
	if err != nil {
		log.Println(err)
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapCreateSaleModelToDto(sale))
}

// SellCollection godoc
// @Summary  sell or hold auction for collections
// @Tags     sale
// @Accept   json
// @Produce  json
// @Router   /v1/sale/sell-collection [post]
// @Param    message  body  dto.SaleRequest  true  "collection sale request body"
// @Success  200      body  model.Sale
func (s SaleController) SellCollection(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "SaleController[SellCollection]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	if c.Body() == nil {
		return filper.GetBadRequestError(c, "you need to provide body in your request")
	}

	var request dto.SaleRequest
	if err := c.BodyParser(&request); err != nil {
		return filper.GetBadRequestError(c, "invalid body data")
	}

	errRes := validator.Validate(request)
	if len(errRes.Errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	sale, err := s.saleService.CreateNftSale(ctx, mapCreateSaleDtoToModel(request, userId))
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapCreateSaleModelToDto(sale))
}

// CancelSale godoc
// @Summary  cancel any type of sale
// @Tags     sale
// @Accept   json
// @Produce  json
// @Param    id  path  string  true  "sale id that you want to be canceled"
// @Router   /v1/sale/{id} [delete]
// @Success  200  {string}  string  "sale canceled"
func (s SaleController) CancelSale(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "SaleController[CancelSale]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	saleId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid sale id")
	}

	if err := s.saleService.CancelSale(ctx, model.Sale{ID: &saleId, User: usermodel.User{ID: userId}}); err != nil {
		return filper.GetInternalError(c, "")
	}

	return filper.GetSuccessResponse(c, "sale canceled")
}

// GetAllSales godoc
// @Summary  get user sale list
// @Tags     sale
// @Accept   json
// @Produce  json
// @Router   /v1/sale [get]
// @Success  200  {object}  dto.SaleList
func (s SaleController) GetAllSales(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "SaleController[GetAllSales]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	list, err := s.saleService.GetSalesList(ctx, userId)
	if err != nil {
		return filper.GetInternalError(c, "")
	}

	if len(list) == 0 {
		return c.Status(fiber.StatusNoContent).SendString("")
	}

	return c.Status(fiber.StatusOK).JSON(createSalesListDtoFromModel(list))
}

// GetSale godoc
// @Summary  sell or hld auction for collections
// @Tags     sale
// @Accept   json
// @Produce  json
// @Param    id  path  string  true  "sale id that you want to retrieve"
// @Router   /v1/sale/{id} [get]
// @Success  200  {object}  dto.Sale
func (s SaleController) GetSale(c *fiber.Ctx) error {
	span, ctx := jtrace.T().SpanFromContext(c.Context(), "SaleController[GetSale]")
	defer span.Finish()

	if c.Locals("user_id") == nil {
		return filper.GetInternalError(c, "")
	}
	userId := c.Locals("user_id").(uuid.UUID)

	saleId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return filper.GetBadRequestError(c, "invalid sale id")
	}

	sale, err := s.saleService.GetSale(ctx, model.Sale{ID: &saleId, User: usermodel.User{ID: userId}})
	if err != nil {
		if errors.Is(err, apperrors.ErrSaleNotFound) {
			return filper.GetNotFoundError(c, "sale not found")
		}
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusOK).JSON(mapSaleModelToDto(sale))
}
