package sale

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"nft/contract"
	"nft/infra/jtrace"
	dto "nft/internal/sale/dto"
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
		return filper.GetInternalError(c, "")
	}

	return c.Status(fiber.StatusCreated).JSON(mapCreateSaleModelToDto(sale))
}

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

func (s SaleController) CancelSale(c *fiber.Ctx) error {
	span, _ := jtrace.T().SpanFromContext(c.Context(), "SaleController[CancelSale]")
	defer span.Finish()
	return nil
}

func (s SaleController) GetAllSales(c *fiber.Ctx) error {
	span, _ := jtrace.T().SpanFromContext(c.Context(), "SaleController[GetAllSales]")
	defer span.Finish()
	return nil
}

func (s SaleController) GetSale(c *fiber.Ctx) error {
	span, _ := jtrace.T().SpanFromContext(c.Context(), "SaleController[GetSale]")
	defer span.Finish()
	return nil
}
